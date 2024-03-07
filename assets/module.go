package assets

import (
	"os"
	"path/filepath"

	"github.com/dangduoc08/ecommerce-api/assets/controllers"
	"github.com/dangduoc08/ecommerce-api/assets/providers"
	"github.com/dangduoc08/gogo/core"
)

var Module = func() *core.Module {
	currentDir, _ := os.Getwd()
	publicPath := filepath.Join(currentDir, "publics")

	mod := core.ModuleBuilder().
		Controllers(
			controllers.REST{
				PublicPath: publicPath,
			},
		).
		Providers(
			providers.HandleAsset{
				PublicPath: publicPath,
			},
		).
		Build()

	mod.OnInit = func() {
		if _, err := os.Stat(publicPath); os.IsNotExist(err) {
			err := os.MkdirAll(publicPath, os.ModePerm)
			if err != nil {
				return
			}
		}
	}

	return mod
}()
