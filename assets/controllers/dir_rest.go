package controllers

import (
	"os"

	"github.com/dangduoc08/ecommerce-api/assets/dtos"
	"github.com/dangduoc08/ecommerce-api/assets/models"
	"github.com/dangduoc08/ecommerce-api/assets/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type DirREST struct {
	common.Guard
	common.REST
	common.Logger
	providers.HandleAsset
	PublicPath string
}

func (self DirREST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("assets")

	self.BindGuard(
		shared.AuthGuard{},
	)

	return self
}

func (self DirREST) READ(
	c gooh.Context,
	queryDTO dtos.READ_Query,
) []*models.Asset {
	ls, err := self.List(self.PublicPath + queryDTO.Dir)

	if err != nil {
		self.Debug(
			"DirREST.READ.List",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)

		return []*models.Asset{}
	}

	return ls
}

func (self DirREST) DELETE(
	c gooh.Context,
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

func (self DirREST) MODIFY(
	c gooh.Context,
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

func (self DirREST) CREATE_dirs(
	c gooh.Context,
	bodyDTO dtos.CREATE_dirs_Body,
) gooh.Map {
	dir, err := self.Mkdir(self.PublicPath, bodyDTO.Data.Name)

	if err != nil {
		self.Debug(
			"DirREST.CREATE_dirs.Mkdir",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return gooh.Map{
		"dir": dir,
	}
}
