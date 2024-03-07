package controllers

import (
	"os"
	"strings"

	"github.com/dangduoc08/ecommerce-api/assets/dtos"
	"github.com/dangduoc08/ecommerce-api/assets/models"
	"github.com/dangduoc08/ecommerce-api/assets/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type REST struct {
	common.Guard
	common.REST
	common.Logger
	providers.HandleAsset
	PublicPath string
}

func (instance REST) NewController() core.Controller {
	instance.
		Prefix("v1").
		Prefix(
			"assets",
			instance.READ,
			instance.CREATE_dirs,
			instance.CREATE_files,
			instance.MODIFY_dirs,
			instance.DELETE,
		).
		Prefix(
			"publics",
			instance.SERVE_ANY,
		)

	instance.BindGuard(
		shared.AuthGuard{},
	)

	return instance
}

func (instance REST) READ(
	ctx gogo.Context,
	queryDTO dtos.READ_Query,
) []*models.Asset {
	ls, err := instance.List(instance.PublicPath + queryDTO.Dir)

	if err != nil {
		instance.Debug(
			"READ.List",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)

		return []*models.Asset{}
	}

	return ls
}

func (instance REST) CREATE_dirs(
	ctx gogo.Context,
	bodyDTO dtos.CREATE_dirs_Body,
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

func (instance REST) CREATE_files(
	fileBodyDTO dtos.CREATE_files_Body,
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

func (instance REST) MODIFY_dirs(
	ctx gogo.Context,
	bodyDTO dtos.MODIFY_Body,
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

func (instance REST) DELETE(
	ctx gogo.Context,
	queryDTO dtos.DELETE_Query,
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

func (instance REST) SERVE_ANY() string {
	return instance.PublicPath
}
