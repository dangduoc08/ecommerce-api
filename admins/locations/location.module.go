package locations

import (
	"github.com/dangduoc08/gogo/core"
)

var LocationModule = core.ModuleBuilder().
	Providers(
		LocationProvider{},
	).
	Build()
