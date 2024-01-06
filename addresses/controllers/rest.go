package controllers

import (
	"github.com/dangduoc08/ecommerce-api/addresses/dtos"
	"github.com/dangduoc08/ecommerce-api/addresses/interceptors"
	"github.com/dangduoc08/ecommerce-api/addresses/models"
	"github.com/dangduoc08/ecommerce-api/addresses/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type REST struct {
	common.Guard
	common.REST
	common.Interceptor
	common.Logger
	providers.DBHandler
}

func (self REST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("addresses")

	self.
		BindGuard(
			shared.AuthGuard{},
		)

	self.
		BindInterceptor(
			interceptors.AddressLocationTransformation{},
		)

	return self
}

func (self REST) READ_BY_id(
	c gooh.Context,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	paramDTO dtos.READ_BY_id_Param,
) *models.Address {
	address, err := self.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: accessTokenPayloadDTO.StoreID,
	})

	if err != nil {
		self.Debug(
			"READ_BY_id.FindOneBy",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return nil
	}

	return address
}

func (self REST) CREATE(
	c gooh.Context,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	dto dtos.CREATE_Body,
) *models.Address {
	dataCreation := &providers.Creation{
		StoreID:    accessTokenPayloadDTO.StoreID,
		StreetName: dto.Data.StreetName,
		LocationID: &dto.Data.LocationID,
	}
	if dto.Data.LocationID == 0 {
		dataCreation.LocationID = nil
	}

	addresses, err := self.CreateOne(dataCreation)

	if err != nil {
		self.Debug(
			"CREATE.CreateOne",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.UnprocessableEntityException(err.Error()))
	}

	return addresses
}

func (self REST) UPDATE_BY_id(
	c gooh.Context,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	paramDTO dtos.UPDATE_BY_id_Param,
	bodyDTO dtos.UPDATE_BY_id_Body,
) *models.Address {
	address, err := self.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: accessTokenPayloadDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	dataUpdate := &providers.Update{
		StoreID:    accessTokenPayloadDTO.StoreID,
		StreetName: bodyDTO.Data.StreetName,
		LocationID: &bodyDTO.Data.LocationID,
	}
	if bodyDTO.Data.LocationID == 0 {
		dataUpdate.LocationID = nil
	}

	address, err = self.UpdateByID(paramDTO.ID, dataUpdate)
	if err != nil {
		self.Debug(
			"UPDATE_BY_id.CreateOne",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.UnprocessableEntityException(err.Error()))
	}

	return address
}

func (self REST) DELETE_BY_id(
	c gooh.Context,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	paramDTO dtos.DELETE_BY_id_Param,
) gooh.Map {
	_, err := self.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: accessTokenPayloadDTO.StoreID,
	})

	if err != nil {
		self.Debug(
			"DELETE_BY_id.FindOneBy",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return gooh.Map{
			"deleted": false,
		}
	}

	if err = self.DeleteByID(paramDTO.ID); err != nil {
		self.Debug(
			"DELETE_BY_id.DeleteByID",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return gooh.Map{
			"deleted": false,
		}
	}

	return gooh.Map{
		"deleted": true,
	}
}
