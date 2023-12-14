package groups

import (
	"net/http"
	"sort"

	"github.com/dangduoc08/ecommerce-api/globals"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
)

type GroupController struct {
	common.REST
	common.Guard
	GroupProvider  GroupProvider
	Logger         common.Logger
	MethodToAction map[string]string
}

func (self GroupController) NewController() core.Controller {
	self.Prefix("v1").Prefix("groups")
	self.BindGuard(globals.AccessAPIGuard{})

	self.MethodToAction = map[string]string{
		http.MethodGet:    "read",
		http.MethodPost:   "create",
		http.MethodPut:    "update",
		http.MethodPatch:  "modify",
		http.MethodDelete: "delete",
	}

	return self
}

func (self GroupController) READ() []*Group {
	groups, err := self.GroupProvider.FindManyBy()

	if err != nil {
		self.Logger.Debug(
			"GroupProvider.FindManyBy",
			"error", err.Error(),
			"groups", groups,
		)
		return []*Group{}
	}

	return groups
}

func (self GroupController) CREATE(bodyDTO CREATE_Body_DTO) *Group {
	self.GroupProvider.CreateOne(&GroupCreation{
		Name:        bodyDTO.Data.Name,
		Permissions: bodyDTO.Data.Permissions,
	})
	return nil
}

func (self GroupController) READ_permissions() any {
	permissions := map[string][]string{}

	for _, restConfiguration := range self.GetConfigurations() {
		permission := permissions[self.MethodToAction[restConfiguration.Method]]
		permission = append(permission, restConfiguration.Method+restConfiguration.Route)
		sort.Sort(sort.StringSlice(permission))
		permissions[self.MethodToAction[restConfiguration.Method]] = permission
	}

	return permissions
}
