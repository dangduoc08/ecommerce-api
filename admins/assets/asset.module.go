package assets

import (
	"os"
	"path/filepath"

	"github.com/dangduoc08/ecommerce-api/admins/assets/commons"
	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/gogo/core"
)

var AssetModule = func() *core.Module {
	currentDir, _ := os.Getwd()
	publicPath := filepath.Join(currentDir, constants.PUBLIC_DIR)

	module := core.ModuleBuilder().
		Controllers(
			AssetController{
				PublicPath: publicPath,
			},
		).
		Providers(
			AssetProvider{
				PublicPath: publicPath,
			},
			commons.CommonProvider{},
		).
		Build()

	module.
		Prefix("assets")

	module.OnInit = func() {
		if _, err := os.Stat(publicPath); os.IsNotExist(err) {
			err := os.MkdirAll(publicPath, os.ModePerm)
			if err != nil {
				return
			}
		}
	}

	return module
}()
