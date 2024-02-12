package providers

import (
	dbProviders "github.com/dangduoc08/ecommerce-api/dbs/providers"
	"github.com/dangduoc08/ecommerce-api/stores/models"
	"github.com/dangduoc08/gooh/core"
	"gorm.io/gorm/clause"
)

type DBHandler struct {
	dbProviders.DB
}

func (instance DBHandler) NewProvider() core.Provider {
	return instance
}

func (instance DBHandler) FindByID(
	id uint,
) (*models.Store, error) {
	storeRec := &models.Store{ID: id}

	if err := instance.Take(storeRec).Error; err != nil {
		return nil, err
	}

	return storeRec, nil
}

func (instance DBHandler) UpdateByID(id uint, data *Update) (*models.Store, error) {
	storeRec := &models.Store{
		ID:          id,
		Name:        data.Name,
		Description: data.Description,
		Phone:       data.Phone,
		Email:       data.Email,
	}

	if err := instance.
		Clauses(clause.Returning{}).
		Updates(&storeRec).
		Error; err != nil {
		return nil, err
	}

	return storeRec, nil
}
