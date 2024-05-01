package stores

import (
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/gogo/core"
	"gorm.io/gorm/clause"
)

type StoreProvider struct {
	dbs.DBProvider
}

func (instance StoreProvider) NewProvider() core.Provider {
	return instance
}

func (instance StoreProvider) FindByID(
	id uint,
) (*StoreModel, error) {
	storeRec := &StoreModel{ID: id}

	if err := instance.Take(storeRec).Error; err != nil {
		return nil, err
	}

	return storeRec, nil
}

func (instance StoreProvider) UpdateByID(id uint, data *Update) (*StoreModel, error) {
	storeRec := &StoreModel{
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
