package stores

import (
	"github.com/dangduoc08/gooh/core"
)

var StoreModule = core.ModuleBuilder().
	Controllers(
		StoreController{},
	).
	Providers(
		StoreProvider{},
	).
	Exports(
		StoreProvider{},
	).
	Build()
