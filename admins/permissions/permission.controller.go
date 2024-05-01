package permissions

import (
	"net/http"
	"sort"
	"strings"

	"github.com/dangduoc08/ecommerce-api/constants"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
)

type PermissionController struct {
	common.REST
	common.Guard
	MethodToAction map[string]string
	Blacklist      []string
}

func (instance PermissionController) NewController() core.Controller {
	instance.
		BindGuard(sharedLayers.AuthGuard{})

	instance.MethodToAction = map[string]string{
		http.MethodGet:    "read",
		http.MethodPost:   "create",
		http.MethodPut:    "update",
		http.MethodPatch:  "modify",
		http.MethodDelete: "delete",
	}

	instance.Blacklist = []string{
		"POST/admins/auths/sessions",
		"POST/admins/auths/tokens",
	}

	return instance
}

func (instance PermissionController) READ_VERSION_1() any {
	permissions := map[string][]string{}

	for _, restConfiguration := range instance.GetConfigurations() {
		if strings.Contains(restConfiguration.Route, constants.ADMIN_PATH) {
			pattern := restConfiguration.Method + restConfiguration.Route
			if utils.ArrIncludes(instance.Blacklist, pattern) {
				continue
			}
			perm := permissions[instance.MethodToAction[restConfiguration.Method]]
			perm = append(perm, pattern)
			sort.Strings(sort.StringSlice(perm))
			permissions[instance.MethodToAction[restConfiguration.Method]] = perm
		}
	}

	return permissions
}
