package providers

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/categories/models"
	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh/core"
)

type DBValidation struct {
	dbProviders.DB
}

func (self DBValidation) NewProvider() core.Provider {
	return self
}

func (self DBValidation) CheckParentCategories(parentCategoryIDs []uint) ([]*models.Category, error) {
	categoryRecs := []*models.Category{}

	if err := self.Where("id IN ?", parentCategoryIDs).Find(&categoryRecs).Error; err != nil {
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
