package providers

import (
	"errors"
	"fmt"

	"github.com/dangduoc08/ecommerce-api/categories/models"
	dbProviders "github.com/dangduoc08/ecommerce-api/dbs/providers"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo/core"
)

type DBValidation struct {
	dbProviders.DB
}

func (instance DBValidation) NewProvider() core.Provider {
	return instance
}

func (instance DBValidation) CheckParentCategories(parentCategoryIDs []uint) ([]*models.Category, error) {
	categoryRecs := []*models.Category{}

	if err := instance.Where("id IN ?", parentCategoryIDs).Find(&categoryRecs).Error; err != nil {
		return nil, err
	}

	storedCategoryIDs := utils.ArrMap(categoryRecs, func(gr *models.Category, index int) uint {
		return gr.ID
	})

	for index, id := range parentCategoryIDs {
		if !utils.ArrIncludes(storedCategoryIDs, id) {
			return nil, fmt.Errorf("%v[%v]=%v doesn't exist", "category_ids", index, id)
		}
	}

	return categoryRecs, nil
}

func (instance DBValidation) CheckCategoryRelationship(categoryID uint, parentCategoryIDs []uint) error {
	for _, parentCategoryID := range parentCategoryIDs {
		if categoryID == parentCategoryID {
			return errors.New(fmt.Sprintf("Circular relationship are not allowed"))
		}

		categoryCategoryRecs := []*models.CategoryCategory{}
		instance.Raw("SELECT * FROM categories_categories WHERE category_id = ?", parentCategoryID).Scan(&categoryCategoryRecs)

		for _, categoryCategoryRec := range categoryCategoryRecs {
			if categoryID == categoryCategoryRec.ParentCategoryID {
				return errors.New(fmt.Sprintf("Circular relationship are not allowed"))
			}

			if categoryCategoryRec.CategoryID != 0 && categoryCategoryRec.ParentCategoryID != 0 {
				return instance.CheckCategoryRelationship(categoryID, []uint{categoryCategoryRec.ParentCategoryID})
			}
		}
	}

	return nil
}
