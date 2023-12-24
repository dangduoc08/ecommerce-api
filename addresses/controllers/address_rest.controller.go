package controllers

import (
	"github.com/dangduoc08/ecommerce-api/addresses/interceptors"
	"github.com/dangduoc08/ecommerce-api/addresses/models"
	"github.com/dangduoc08/ecommerce-api/addresses/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type AddressREST struct {
	common.Guard
	common.REST
	common.Interceptor
	AddressDB providers.AddressDB
	Logger    common.Logger
}

func (self AddressREST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("addresses")

	self.
		BindGuard(
			shared.AuthGuard{},
			self.READ_BY_id,
			self.CREATE,
			self.UPDATE_BY_id,
			self.DELETE_BY_id,
		)

	self.
		BindInterceptor(
			interceptors.AddressLocationTransformation{},
		)

	return self
}

func (self AddressREST) READ(
	c gooh.Context,
	queryDTO models.READ_Query,
) []*models.Address {
	addresses, err := self.AddressDB.FindManyBy(&providers.AddressQuery{
		StoreID: queryDTO.StoreID,
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
	})

	if err != nil {
		self.Logger.Debug(
			"AddressREST.READ.AddressDB.FindManyBy",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return []*models.Address{}
	}

	return addresses
}

func (self AddressREST) READ_BY_id(
	c gooh.Context,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	paramDTO models.READ_BY_id_Param,
) *models.Address {
	address, err := self.AddressDB.FindOneBy(&providers.AddressQuery{
		ID:      paramDTO.ID,
		StoreID: accessTokenPayloadDTO.StoreID,
	})

	if err != nil {
		self.Logger.Debug(
			"AddressREST.READ_BY_id.AddressDB.FindOneBy",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return nil
	}

	return address
}

func (self AddressREST) CREATE(
	c gooh.Context,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	dto models.CREATE_Body,
) *models.Address {
	dataCreation := &providers.AddressCreation{
		StoreID:    accessTokenPayloadDTO.StoreID,
		StreetName: dto.Data.StreetName,
		LocationID: &dto.Data.LocationID,
	}
	if dto.Data.LocationID == 0 {
		dataCreation.LocationID = nil
	}

	addresses, err := self.AddressDB.CreateOne(dataCreation)

	if err != nil {
		self.Logger.Debug(
			"AddressREST.CREATE.AddressDB.CreateOne",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.UnprocessableEntityException(err.Error()))
	}

	return addresses
}

func (self AddressREST) UPDATE_BY_id(
	c gooh.Context,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	paramDTO models.UPDATE_BY_id_Param,
	bodyDTO models.UPDATE_BY_id_Body,
) *models.Address {
	address, err := self.AddressDB.FindOneBy(&providers.AddressQuery{
		ID:      paramDTO.ID,
		StoreID: accessTokenPayloadDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	dataUpdate := &providers.AddressUpdate{
		StoreID:    accessTokenPayloadDTO.StoreID,
		StreetName: bodyDTO.Data.StreetName,
		LocationID: &bodyDTO.Data.LocationID,
	}
	if bodyDTO.Data.LocationID == 0 {
		dataUpdate.LocationID = nil
	}

	address, err = self.AddressDB.UpdateByID(paramDTO.ID, dataUpdate)
	if err != nil {
		self.Logger.Debug(
			"AddressREST.UPDATE_BY_id.AddressDB.CreateOne",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.UnprocessableEntityException(err.Error()))
	}

	return address
}

func (self AddressREST) DELETE_BY_id(
	c gooh.Context,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	paramDTO models.DELETE_BY_id_Param,
) gooh.Map {
	_, err := self.AddressDB.FindOneBy(&providers.AddressQuery{
		ID:      paramDTO.ID,
		StoreID: accessTokenPayloadDTO.StoreID,
	})

	if err != nil {
		self.Logger.Debug(
			"AddressREST.DELETE_BY_id.AddressDB.FindOneBy",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return gooh.Map{
			"deleted": false,
		}
	}

	if err = self.AddressDB.DeleteByID(paramDTO.ID); err != nil {
		self.Logger.Debug(
			"AddressREST.DELETE_BY_id.AddressDB.DeleteByID",
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
