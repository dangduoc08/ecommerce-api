package users

import (
	authProviders "github.com/dangduoc08/ecommerce-api/auths/providers"
	"github.com/dangduoc08/ecommerce-api/groups"
	"github.com/dangduoc08/ecommerce-api/users/controllers"
	"github.com/dangduoc08/ecommerce-api/users/providers"
	"github.com/dangduoc08/gooh/core"
)

var Module = core.ModuleBuilder().
	Imports(
		groups.Module,
	).
	Controllers(
		controllers.UserREST{},
	).
	Providers(
		providers.UserDB{},
		authProviders.AuthHandler{},
	).
	Build()
