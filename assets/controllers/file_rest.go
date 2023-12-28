package controllers

import (
	"strings"

	"github.com/dangduoc08/ecommerce-api/assets/dtos"
	"github.com/dangduoc08/ecommerce-api/assets/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
)

type FileREST struct {
	common.Guard
	common.REST
	common.Logger
	providers.HandleAsset
	PublicPath string
}

func (self FileREST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("assets")

	self.BindGuard(
		shared.AuthGuard{},
	)

	return self
}

func (self FileREST) CREATE_files(
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
