package providers

import (
	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	"github.com/dangduoc08/ecommerce-api/stores/models"
	"github.com/dangduoc08/gooh/core"
	"gorm.io/gorm/clause"
)

type DBHandler struct {
	dbProviders.DB
}

func (self DBHandler) NewProvider() core.Provider {
	return self
}

func (self DBHandler) FindByID(
	id uint,
) (*models.Store, error) {
	storeRec := &models.Store{ID: id}

	if err := self.Take(storeRec).Error; err != nil {
		return nil, err
	}

	return storeRec, nil
}

func (self DBHandler) UpdateByID(id uint, data *Update) (*models.Store, error) {
	storeRec := &models.Store{
		ID:          id,
		Name:        data.Name,
		Description: data.Description,
		Phone:       data.Phone,
		Email:       data.Email,
	}

	if err := self.
		Clauses(clause.Returning{}).
		Updates(&storeRec).
		Error; err != nil {
		return nil, err
	}

	return storeRec, nil
}
