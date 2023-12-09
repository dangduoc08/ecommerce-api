package addresses

import (
	"path"
	"reflect"

	"github.com/dangduoc08/ecommerce-api/db"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
)

type AddressProvider struct {
	DBProvider    db.DBProvider
	ConfigService config.ConfigService
}

func (self AddressProvider) NewProvider() core.Provider {
	return self
}

func (self AddressProvider) GetModelName() string {
	return path.Base(reflect.TypeOf(self).PkgPath())
}
