package user

import (
	"github.com/dangduoc08/ecommerce-api/database"
	"github.com/dangduoc08/gooh/core"
)

var Module = func() *core.Module {
	return core.ModuleBuilder().
		Imports(
			database.Module,
		).
		Controllers(
			Controller{},
		).
		Providers(
			Provider{},
		).
		Build()

}()
