package providers

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/categories/models"
	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	"github.com/dangduoc08/gooh/core"
	"gorm.io/gorm/clause"
)

type DBHandler struct {
	dbProviders.DB
	DBValidation
}

func (self DBHandler) NewProvider() core.Provider {
	return self
}

func (self DBHandler) FindByID(id uint) (*models.Category, error) {
	categoryRec := &models.Category{
		ID: id,
	}

	if err := self.
		Preload("ParentCategories").
		First(categoryRec).
		Error; err != nil {
		return nil, err
	}

	return categoryRec, nil
}

func (self DBHandler) FindOneBy(query *Query) (*models.Category, error) {
	categoryRec := &models.Category{}
	categoryQueries := map[string]any{}

	if query.ID != 0 {
		categoryQueries["id"] = query.ID
	}

	if query.StoreID != 0 {
		categoryQueries["store_id"] = query.StoreID
	}

	if err := self.
		Where(categoryQueries).
		First(categoryRec).
		Error; err != nil {
		return nil, err
	}

	return categoryRec, nil
}

func (self DBHandler) FindManyBy(query *Query) ([]*models.Category, error) {
	categoryRecs := []*models.Category{}
	categoryQueries := map[string]any{}
	tx := self.DB.DB

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

func (self DBHandler) CreateOne(data *Creation) (*models.Category, error) {
	categoryRec := &models.Category{
		Name:            data.Name,
		Description:     data.Description,
		StoreID:         data.StoreID,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            data.Slug,
		Status:          data.Status,
	}

	if parentCategories, err := self.CheckParentCategories(data.ParentCategoryIDs); err != nil {
		return nil, err
	} else {
		categoryRec.ParentCategories = parentCategories
	}

	if err := self.
		Create(&categoryRec).
		Error; err != nil {
		return nil, err
	}

	return categoryRec, nil
}

func (self DBHandler) UpdateByID(id uint, data *Update) (*models.Category, error) {
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
	if parentCategories, err := self.CheckParentCategories(data.ParentCategoryIDs); err != nil {
		return nil, err
	} else {
		categoryRec.ParentCategories = parentCategories
	}

	if err := self.
		Clauses(clause.Returning{}).
		Updates(&categoryRec).
		Error; err != nil {
		return nil, err
	}

	return categoryRec, nil
}

func (self DBHandler) FindManyAsMenu(query *Query) (*[]map[string]any, error) {
	menu := &[]map[string]any{}
	args := []any{
		constants.CATEGORY_ENABLED,
		constants.CATEGORY_ENABLED,
		query.StoreID,
	}

	base := `
		SELECT
			categories.*,
			json_agg(to_json(child)) AS child_categories
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

	if err := self.
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
