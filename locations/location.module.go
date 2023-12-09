package locations

import (
	"github.com/dangduoc08/gooh/core"
)

var LocationModule = core.ModuleBuilder().
	Controllers(
		LocationController{},
	).
	Providers(
		LocationProvider{},
	).
	Exports(
		LocationProvider{},
	).
	Build()
