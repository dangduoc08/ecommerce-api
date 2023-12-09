package addresses

import (
	"github.com/dangduoc08/ecommerce-api/globals"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
)

type AddressController struct {
	common.Guard
	common.REST
	common.Interceptor
	ConfigService   config.ConfigService
	AddressProvider AddressProvider
	Logger          common.Logger
}

func (self AddressController) NewController() core.Controller {
	self.Prefix("v1").Prefix("addresses")
	self.BindGuard(globals.AccessAPIGuard{}, self.CREATE, self.UPDATE_BY_id)
	self.BindInterceptor(AddressInterceptor{})

	return self
}

func (self AddressController) READ(ctx gooh.Context, queryDTO READ_Query) []Address {
	addresses, err := self.AddressProvider.FindManyBy(&AddressQuery{
		StoreID: queryDTO.StoreID,
	})

	if err != nil {
		self.Logger.Debug(
			"AddressProvider.FindManyBy",
			"error", err.Error(),
			"addresses", addresses,
			"X-Request-ID", ctx.GetID(),
		)
		return []Address{}
	}

	return addresses
}

func (self AddressController) CREATE(
	ctx gooh.Context,
	accessTokenDTO globals.AccessTokenDTO,
	dto CREATE_Body_DTO,
) *Address {
	dataCreation := &AddressCreation{
		StoreID:    accessTokenDTO.StoreID,
		StreetName: dto.Data.StreetName,
		LocationID: &dto.Data.LocationID,
	}
	if dto.Data.LocationID == 0 {
		dataCreation.LocationID = nil
	}

	addresses, err := self.AddressProvider.CreateOne(dataCreation)

	if err != nil {
		self.Logger.Debug(
			"AddressProvider.CreateOne",
			"error", err.Error(),
			"addresses", addresses,
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.UnprocessableEntityException(err.Error()))
	}

	return addresses
}

func (self AddressController) UPDATE_BY_id(
	ctx gooh.Context,
	accessTokenDTO globals.AccessTokenDTO,
	paramDTO UPDATE_BY_ID_Param_DTO,
	bodyDTO UPDATE_BY_ID_Body_DTO,
) *Address {

	addresses, err := self.AddressProvider.FindManyBy(&AddressQuery{
		ID:      paramDTO.ID,
		StoreID: accessTokenDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	if len(addresses) == 0 {
		panic(exception.NotFoundException("Field: id, Error: not exists"))
	}

	dataUpdate := &AddressUpdate{
		StoreID:    accessTokenDTO.StoreID,
		StreetName: bodyDTO.Data.StreetName,
		LocationID: &bodyDTO.Data.LocationID,
	}
	if bodyDTO.Data.LocationID == 0 {
		dataUpdate.LocationID = nil
	}

	address, err := self.AddressProvider.UpdateByID(paramDTO.ID, dataUpdate)
	if err != nil {
		self.Logger.Debug(
			"AddressProvider.UpdateByID",
			"error", err.Error(),
			"address", address,
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.UnprocessableEntityException(err.Error()))
	}

	return address
}
