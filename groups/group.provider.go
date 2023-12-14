package groups

import (
	"path"
	"reflect"

	"github.com/dangduoc08/ecommerce-api/db"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
)

type GroupProvider struct {
	DBProvider    db.DBProvider
	ConfigService config.ConfigService
	Logger        common.Logger
}

func (self GroupProvider) NewProvider() core.Provider {
	return self
}

func (self GroupProvider) GetModelName() string {
	return path.Base(reflect.TypeOf(self).PkgPath())
}

func (self GroupProvider) FindManyBy() ([]*Group, error) {
	var groups []*Group

	resp := self.DBProvider.DB.Find(&groups)

	return groups, resp.Error
}

func (self GroupProvider) CreateOne(data *GroupCreation) (*Group, error) {
	group := &Group{
		Name:        data.Name,
		Permissions: data.Permissions,
	}

	resp := self.DBProvider.DB.Create(group)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return group, nil
}
