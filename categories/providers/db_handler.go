package providers

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/categories/models"
	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/dbs/providers"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo/core"
	"gorm.io/gorm/clause"
)

type DBHandler struct {
	dbProviders.DB
	DBValidation
}

func (instance DBHandler) NewProvider() core.Provider {
	return instance
}

func (instance DBHandler) FindByID(id uint) (*models.Category, error) {
	categoryRec := &models.Category{
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

func (instance DBHandler) FindOneBy(query *Query) (*models.Category, error) {
	categoryRec := &models.Category{}
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

func (instance DBHandler) FindManyBy(query *Query) ([]*models.Category, error) {
	categoryRecs := []*models.Category{}
	categoryQueries := map[string]any{}
	tx := instance.DB.DB

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
		return []*models.Category{}, err
	}

	return categoryRecs, nil
}

func (instance DBHandler) CreateOne(data *Creation) (*models.Category, error) {
	categoryRec := &models.Category{
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

func (instance DBHandler) UpdateByID(id uint, data *Update) (*models.Category, error) {
	categoryRec := &models.Category{
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

	tx := instance.DB.Begin()
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

func (instance DBHandler) FindManyAsMenu(query *Query) (*[]map[string]any, error) {
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

func (instance DBHandler) FindCurrentParentCategoryIDs(categoryID uint) []uint {
	categoryCategoryRecs := []*models.CategoryCategory{}
	instance.Raw("SELECT * FROM categories_categories WHERE category_id = ?", categoryID).Scan(&categoryCategoryRecs)

	return utils.ArrMap(categoryCategoryRecs, func(categoryCategoryRec *models.CategoryCategory, i int) uint {
		return categoryCategoryRec.ParentCategoryID
	})
}
