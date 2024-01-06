package controllers

import (
	"net/http"
	"sort"

	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
)

type REST struct {
	common.REST
	common.Guard
	MethodToAction map[string]string
	Blacklist      []string
}

func (self REST) NewController() core.Controller {
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

	self.Blacklist = []string{
		"POST/v1/auths/sessions",
		"POST/v1/auths/tokens",
		"GET/v1/stores/{id}",
		"GET/v1/stores/{id}/addresses",
		"GET/v1/stores/{id}/categories",
	}

	return self
}

func (self REST) READ() any {
	permissions := map[string][]string{}

	for _, restConfiguration := range self.GetConfigurations() {
		pattern := restConfiguration.Method + restConfiguration.Route
		if utils.ArrIncludes(self.Blacklist, pattern) {
			continue
		}
		perm := permissions[self.MethodToAction[restConfiguration.Method]]
		perm = append(perm, pattern)
		sort.Sort(sort.StringSlice(perm))
		permissions[self.MethodToAction[restConfiguration.Method]] = perm
	}

	return permissions
}
