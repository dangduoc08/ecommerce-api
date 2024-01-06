package controllers

import (
	"os"
	"strings"

	"github.com/dangduoc08/ecommerce-api/assets/dtos"
	"github.com/dangduoc08/ecommerce-api/assets/models"
	"github.com/dangduoc08/ecommerce-api/assets/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type REST struct {
	common.Guard
	common.REST
	common.Logger
	providers.HandleAsset
	PublicPath string
}

func (self REST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("assets")

	self.BindGuard(
		shared.AuthGuard{},
	)

	return self
}

func (self REST) READ(
	ctx gooh.Context,
	queryDTO dtos.READ_Query,
) []*models.Asset {
	ls, err := self.List(self.PublicPath + queryDTO.Dir)

	if err != nil {
		self.Debug(
			"READ.List",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)

		return []*models.Asset{}
	}

	return ls
}

func (self REST) CREATE_dirs(
	ctx gooh.Context,
	bodyDTO dtos.CREATE_dirs_Body,
) gooh.Map {
	dir, err := self.Mkdir(self.PublicPath, bodyDTO.Data.Dir)

	if err != nil {
		self.Debug(
			"CREATE_dirs.Mkdir",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return gooh.Map{
		"dir": dir,
	}
}

func (self REST) CREATE_files(
	fileBodyDTO dtos.CREATE_files_Body,
) gooh.Map {
	if fileBodyDTO.File.Dest == "" {
		return gooh.Map{
			"uploaded": false,
			"dir":      "",
		}
	}

	dir := strings.Replace(fileBodyDTO.File.Dest, self.PublicPath, "", 1)
	if dir != "" {
		dir = dir[1:]
	}

	return gooh.Map{
		"uploaded": true,
		"dir":      dir,
	}
}

func (self REST) MODIFY_names(
	ctx gooh.Context,
	bodyDTO dtos.MODIFY_Body,
) gooh.Map {
	if err := os.Rename(self.PublicPath+bodyDTO.Data.OldDir, self.PublicPath+bodyDTO.Data.NewDir); err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	dir := bodyDTO.Data.NewDir
	if dir != "" {
		dir = dir[1:]
	}
	return gooh.Map{
		"dir": dir,
	}
}

func (self REST) DELETE(
	ctx gooh.Context,
	queryDTO dtos.DELETE_Query,
) gooh.Map {
	for _, dir := range queryDTO.Dirs {
		if dir != "/" {
			if err := os.RemoveAll(self.PublicPath + dir); err != nil {
				panic(exception.InternalServerErrorException(err.Error()))
			}
		}
	}

	return gooh.Map{
		"deleted": true,
	}
}
