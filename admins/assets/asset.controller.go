package assets

import (
	"os"
	"strings"

	"github.com/dangduoc08/ecommerce-api/admins/assets/dtos"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type AssetController struct {
	common.Guard
	common.REST
	common.Logger
	AssetProvider
	PublicPath string
}

func (instance AssetController) NewController() core.Controller {
	instance.
		BindGuard(sharedLayers.AuthGuard{})

	return instance
}

func (instance AssetController) READ_VERSION_1(
	ctx gogo.Context,
	queryDTO dtos.READ_Query_DTO,
) []*AssetModel {
	ls, err := instance.List(instance.PublicPath + queryDTO.Dir)

	if err != nil {
		instance.Debug(
			"READ.List",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)

		return []*AssetModel{}
	}

	return ls
}

func (instance AssetController) CREATE_dirs_VERSION_1(
	ctx gogo.Context,
	bodyDTO dtos.CREATE_dirs_Body_DTO,
) gogo.Map {
	dir, err := instance.Mkdir(instance.PublicPath, bodyDTO.Data.Dir)

	if err != nil {
		instance.Debug(
			"CREATE_dirs.Mkdir",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return gogo.Map{
		"dir": dir,
	}
}

func (instance AssetController) CREATE_files_VERSION_1(
	fileBodyDTO dtos.CREATE_files_Body_DTO,
) gogo.Map {
	if fileBodyDTO.File.Dest == "" {
		return gogo.Map{
			"uploaded": false,
			"dir":      "",
		}
	}

	dir := strings.Replace(fileBodyDTO.File.Dest, instance.PublicPath, "", 1)
	if dir != "" {
		dir = dir[1:]
	}

	return gogo.Map{
		"uploaded": true,
		"dir":      dir,
	}
}

func (instance AssetController) MODIFY_dirs_VERSION_1(
	ctx gogo.Context,
	bodyDTO dtos.MODIFY_Body_DTO,
) gogo.Map {
	if err := os.Rename(instance.PublicPath+bodyDTO.Data.OldDir, instance.PublicPath+bodyDTO.Data.NewDir); err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	dir := bodyDTO.Data.NewDir
	if dir != "" {
		dir = dir[1:]
	}
	return gogo.Map{
		"dir": dir,
	}
}

func (instance AssetController) DELETE_VERSION_1(
	ctx gogo.Context,
	queryDTO dtos.DELETE_Query_DTO,
) gogo.Map {
	for _, dir := range queryDTO.Dirs {
		if dir != "/" {
			if err := os.RemoveAll(instance.PublicPath + dir); err != nil {
				panic(exception.InternalServerErrorException(err.Error()))
			}
		}
	}

	return gogo.Map{
		"deleted": true,
	}
}
