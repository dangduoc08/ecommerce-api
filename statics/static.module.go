package statics

import (
	"os"
	"path/filepath"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/gogo/core"
)

var StaticModule = func() *core.Module {
	currentDir, _ := os.Getwd()
	publicPath := filepath.Join(currentDir, constants.PUBLIC_DIR)

	module := core.ModuleBuilder().
		Controllers(
			StaticController{
				PublicPath: publicPath,
			},
		).
		Build()

	module.
		Prefix("statics")

	return module
}()
