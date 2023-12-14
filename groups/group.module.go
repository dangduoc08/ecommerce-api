package groups

import "github.com/dangduoc08/gooh/core"

var GroupModule = core.ModuleBuilder().
	Controllers(
		GroupController{},
	).
	Providers(
		GroupProvider{},
	).
	Exports(
		GroupProvider{},
	).
	Build()
