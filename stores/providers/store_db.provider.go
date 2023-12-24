package providers

import (
	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	"github.com/dangduoc08/ecommerce-api/stores/models"
	"github.com/dangduoc08/gooh/core"
	"gorm.io/gorm/clause"
)

type StoreDB struct {
	DBProvider dbProviders.DB
}

func (self StoreDB) NewProvider() core.Provider {
	return self
}

func (self StoreDB) FindByID(
	id uint,
) (*models.Store, error) {
	storeRec := &models.Store{ID: id}

	if err := self.DBProvider.DB.Take(storeRec).Error; err != nil {
		return nil, err
	}

	return storeRec, nil
}

func (self StoreDB) UpdateByID(id uint, data *StoreUpdate) (*models.Store, error) {
	storeRec := &models.Store{
		ID:          id,
		Name:        data.Name,
		Description: data.Description,
		Phone:       data.Phone,
		Email:       data.Email,
	}

	if err := self.DBProvider.DB.
		Clauses(clause.Returning{}).
		Updates(&storeRec).
		Error; err != nil {
		return nil, err
	}

	return storeRec, nil
}
