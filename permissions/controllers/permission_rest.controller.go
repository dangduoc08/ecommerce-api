package controllers

import (
	"net/http"
	"sort"

	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
)

type PermissionREST struct {
	common.REST
	common.Guard
	MethodToAction map[string]string
}

func (self PermissionREST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("permissions")

	self.
		BindGuard(shared.AuthGuard{})

	self.MethodToAction = map[string]string{
		http.MethodGet:    "read",
		http.MethodPost:   "create",
		http.MethodPut:    "update",
		http.MethodPatch:  "modify",
		http.MethodDelete: "delete",
	}

	return self
}

func (self PermissionREST) READ() any {
	permissions := map[string][]string{}

	for _, restConfiguration := range self.GetConfigurations() {
		permission := permissions[self.MethodToAction[restConfiguration.Method]]
		permission = append(permission, restConfiguration.Method+restConfiguration.Route)
		sort.Sort(sort.StringSlice(permission))
		permissions[self.MethodToAction[restConfiguration.Method]] = permission
	}

	return permissions
}
