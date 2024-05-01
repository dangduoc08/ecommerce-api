package categories

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo/core"
	"gorm.io/gorm/clause"
)

type CategoryProvider struct {
	dbs.DBProvider
}

func (instance CategoryProvider) NewProvider() core.Provider {
	return instance
}

func (instance CategoryProvider) FindByID(id uint) (*CategoryModel, error) {
	categoryRec := &CategoryModel{
		ID: id,
	}

	if err := instance.
		Preload("ParentCategories").
		First(categoryRec).
		Error; err != nil {
		return nil, err
	}

	return categoryRec, nil
}

func (instance CategoryProvider) FindOneBy(query *Query) (*CategoryModel, error) {
	categoryRec := &CategoryModel{}
	categoryQueries := map[string]any{}

	if query.ID != 0 {
		categoryQueries["id"] = query.ID
	}

	if query.StoreID != 0 {
		categoryQueries["store_id"] = query.StoreID
	}

	if err := instance.
		Where(categoryQueries).
		First(categoryRec).
		Error; err != nil {
		return nil, err
	}

	return categoryRec, nil
}

func (instance CategoryProvider) FindManyBy(query *Query) ([]*CategoryModel, error) {
	categoryRecs := []*CategoryModel{}
	categoryQueries := map[string]any{}
	tx := instance.DBProvider.DB

	if query.StoreID != 0 {
		categoryQueries["store_id"] = query.StoreID
	}

	if query.Status != "" {
		categoryQueries["status"] = query.Status
	}

	if query.Sort != "" {
		if query.Order == "" {
			query.Order = constants.ASC
		}
		tx = tx.Order(fmt.Sprintf("%v %v", query.Sort, query.Order))
	}

	if err := tx.
		Limit(query.Limit).
		Offset(query.Offset).
		Where(categoryQueries).
		Preload("ParentCategories").
		Find(&categoryRecs).
		Error; err != nil {
		return []*CategoryModel{}, err
	}

	return categoryRecs, nil
}

func (instance CategoryProvider) CreateOne(data *Creation) (*CategoryModel, error) {
	categoryRec := &CategoryModel{
		Name:            data.Name,
		Description:     data.Description,
		StoreID:         data.StoreID,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            data.Slug,
		Status:          data.Status,
	}

	if parentCategories, err := instance.CheckParentCategories(data.ParentCategoryIDs); err != nil {
		return nil, err
	} else {
		categoryRec.ParentCategories = parentCategories
	}

	if err := instance.
		Create(&categoryRec).
		Error; err != nil {
		return nil, err
	}

	return categoryRec, nil
}

func (instance CategoryProvider) UpdateByID(id uint, data *Update) (*CategoryModel, error) {
	categoryRec := &CategoryModel{
		ID:              id,
		StoreID:         data.StoreID,
		Name:            data.Name,
		Description:     data.Description,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            data.Slug,
		Status:          data.Status,
	}

	if parentCategories, err := instance.CheckParentCategories(data.ParentCategoryIDs); err != nil {
		return nil, err
	} else {
		categoryRec.ParentCategories = parentCategories
	}

	if err := instance.CheckCategoryRelationship(id, data.ParentCategoryIDs); err != nil {
		return nil, err
	}

	tx := instance.DBProvider.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	currentParentCategoryIDs := instance.FindCurrentParentCategoryIDs(id)
	for _, currentParentCategoryID := range currentParentCategoryIDs {
		if !utils.ArrIncludes(data.ParentCategoryIDs, currentParentCategoryID) {
			if err := tx.Exec("DELETE FROM categories_categories WHERE parent_category_id = ?", currentParentCategoryID).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	if err := tx.
		Clauses(clause.Returning{}).
		Updates(&categoryRec).
		Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return categoryRec, tx.Commit().Error
}

func (instance CategoryProvider) FindManyAsMenu(query *Query) (*[]map[string]any, error) {
	menu := &[]map[string]any{}
	args := []any{
		constants.CATEGORY_ENABLED,
		constants.CATEGORY_ENABLED,
		query.StoreID,
	}

	base := `
		SELECT
			categories.*,
			json_agg(to_json(child)) AS sub_categories
		FROM
			categories
			LEFT JOIN categories_categories
				ON categories.id = categories_categories.parent_category_id
			LEFT JOIN categories child
				ON categories_categories.category_id = child.id
				AND child.status = ?
	`

	conditionals := `
		WHERE categories.status = ?
		AND categories.id NOT IN (
			SELECT
				category_id
			FROM
				categories_categories
			WHERE
				category_id IS NOT NULL
		)
		AND categories.store_id = ?
	`

	if query.CategoryID != 0 {
		args = append(args, query.CategoryID)

		conditionals = `
			WHERE categories.status = ?
			AND categories.store_id = ?
			AND categories.id = ?
		`
	}

	if err := instance.
		Raw(
			base+conditionals+"GROUP BY categories.id",
			args...,
		).
		Scan(menu).
		Error; err != nil {
		return &[]map[string]any{}, err
	}

	return menu, nil
}

func (instance CategoryProvider) FindCurrentParentCategoryIDs(categoryID uint) []uint {
	categoryCategoryRecs := []*CategoryCategoryModel{}
	instance.Raw("SELECT * FROM categories_categories WHERE category_id = ?", categoryID).Scan(&categoryCategoryRecs)

	return utils.ArrMap(categoryCategoryRecs, func(categoryCategoryRec *CategoryCategoryModel, i int) uint {
		return categoryCategoryRec.ParentCategoryID
	})
}

func (instance CategoryProvider) CheckParentCategories(parentCategoryIDs []uint) ([]*CategoryModel, error) {
	categoryRecs := []*CategoryModel{}

	if err := instance.Where("id IN ?", parentCategoryIDs).Find(&categoryRecs).Error; err != nil {
		return nil, err
	}

	storedCategoryIDs := utils.ArrMap(categoryRecs, func(gr *CategoryModel, index int) uint {
		return gr.ID
	})

	for index, id := range parentCategoryIDs {
		if !utils.ArrIncludes(storedCategoryIDs, id) {
			return nil, fmt.Errorf("%v[%v]=%v doesn't exist", "category_ids", index, id)
		}
	}

	return categoryRecs, nil
}

func (instance CategoryProvider) CheckCategoryRelationship(categoryID uint, parentCategoryIDs []uint) error {
	for _, parentCategoryID := range parentCategoryIDs {
		if categoryID == parentCategoryID {
			return fmt.Errorf("circular relationship are not allowed")
		}

		categoryCategoryRecs := []*CategoryCategoryModel{}
		instance.Raw("SELECT * FROM categories_categories WHERE category_id = ?", parentCategoryID).Scan(&categoryCategoryRecs)

		for _, categoryCategoryRec := range categoryCategoryRecs {
			if categoryID == categoryCategoryRec.ParentCategoryID {
				return fmt.Errorf("circular relationship are not allowed")
			}

			if categoryCategoryRec.CategoryID != 0 && categoryCategoryRec.ParentCategoryID != 0 {
				return instance.CheckCategoryRelationship(categoryID, []uint{categoryCategoryRec.ParentCategoryID})
			}
		}
	}

	return nil
}
