package users

import (
	authProviders "github.com/dangduoc08/ecommerce-api/auths/providers"
	"github.com/dangduoc08/ecommerce-api/groups"
	"github.com/dangduoc08/ecommerce-api/users/controllers"
	"github.com/dangduoc08/ecommerce-api/users/providers"
	"github.com/dangduoc08/gogo/core"
)

var Module = core.ModuleBuilder().
	Imports(
		groups.Module,
	).
	Controllers(
		controllers.REST{},
	).
	Providers(
		providers.DBHandler{},
		providers.DBValidation{},
		authProviders.Cipher{},
	).
	Build()
