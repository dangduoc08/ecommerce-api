package controllers

import (
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/ecommerce-api/stores/models"
	"github.com/dangduoc08/ecommerce-api/stores/providers"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type StoreREST struct {
	common.REST
	common.Guard
	Logger  common.Logger
	StoreDB providers.StoreDB
}

func (self StoreREST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("stores")

	self.BindGuard(
		shared.AuthGuard{},
		self.UPDATE_BY_id,
	)

	return self
}

func (self StoreREST) READ_BY_id(
	c gooh.Context,
	paramDTO models.READ_BY_id_Param,
) *models.Store {
	store, err := self.StoreDB.FindByID(paramDTO.ID)
	if err != nil {
		self.Logger.Debug(
			"StoreREST.READ_BY_id.StoreDB.FindByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	return store
}

func (self StoreREST) UPDATE_BY_id(
	c gooh.Context,
	paramDTO models.UPDATE_BY_id_Param,
	bodyDTO models.UPDATE_BY_id_Body,
) *models.Store {
	store, err := self.StoreDB.FindByID(paramDTO.ID)
	if err != nil {
		self.Logger.Debug(
			"StoreREST.UPDATE_BY_id.StoreDB.FindByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	store, err = self.StoreDB.UpdateByID(paramDTO.ID, &providers.StoreUpdate{
		Name:        bodyDTO.Data.Name,
		Description: bodyDTO.Data.Description,
		Phone:       bodyDTO.Data.Phone,
		Email:       bodyDTO.Data.Email,
	})

	if err != nil {
		self.Logger.Debug(
			"GroupREST.UPDATE_BY_id.StoreDB.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return store
}
