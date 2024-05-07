package seeds

import (
	"github.com/dangduoc08/ecommerce-api/admins/auths"
	"github.com/dangduoc08/gogo/core"
)

var SeedModule = core.ModuleBuilder().
	Imports(
		auths.AuthModule,
	).
	Providers(
		SeedProvider{},
	).
	Build()
