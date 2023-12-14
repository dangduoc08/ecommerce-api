package stores

import (
	"path"
	"reflect"

	"github.com/dangduoc08/ecommerce-api/db"
	"github.com/dangduoc08/gooh/core"
)

type StoreProvider struct {
	DBProvider db.DBProvider
}

func (self StoreProvider) NewProvider() core.Provider {
	return self
}

func (self StoreProvider) GetModelName() string {
	return path.Base(reflect.TypeOf(self).PkgPath())
}

func (self StoreProvider) FindOneByID(
	id uint,
) (*Store, error) {
	store := &Store{ID: id}
	resp := self.DBProvider.DB.Take(store)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return store, nil
}