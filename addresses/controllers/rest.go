package controllers

import (
	"github.com/dangduoc08/ecommerce-api/addresses/dtos"
	"github.com/dangduoc08/ecommerce-api/addresses/interceptors"
	"github.com/dangduoc08/ecommerce-api/addresses/models"
	"github.com/dangduoc08/ecommerce-api/addresses/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type REST struct {
	common.Guard
	common.REST
	common.Interceptor
	common.Logger
	providers.DBHandler
}

func (instance REST) NewController() core.Controller {
	instance.
		Prefix("v1").
		Prefix("addresses")

	instance.
		BindGuard(shared.AuthGuard{})

	instance.
		BindInterceptor(interceptors.AddressLocationTransformation{})

	return instance
}

func (instance REST) READ_BY_id(
	ctx gogo.Context,
	tokenClaimsDTO shared.TokenClaimsDTO,
	paramDTO dtos.READ_BY_id_Param,
) *models.Address {
	address, err := instance.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		instance.Debug(
			"READ_BY_id.FindOneBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return nil
	}

	return address
}

func (instance REST) CREATE(
	ctx gogo.Context,
	tokenClaimsDTO shared.TokenClaimsDTO,
	dto dtos.CREATE_Body,
) *models.Address {
	dataCreation := &providers.Creation{
		StoreID:    tokenClaimsDTO.StoreID,
		StreetName: &dto.Data.StreetName,
		LocationID: &dto.Data.LocationID,
	}
	if dto.Data.LocationID == 0 {
		dataCreation.LocationID = nil
	}

	addresses, err := instance.CreateOne(dataCreation)

	if err != nil {
		instance.Debug(
			"CREATE.CreateOne",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.UnprocessableEntityException(err.Error()))
	}

	return addresses
}

func (instance REST) UPDATE_BY_id(
	ctx gogo.Context,
	tokenClaimsDTO shared.TokenClaimsDTO,
	paramDTO dtos.UPDATE_BY_id_Param,
	bodyDTO dtos.UPDATE_BY_id_Body,
) *models.Address {
	address, err := instance.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	dataUpdate := &providers.Update{
		StoreID:    tokenClaimsDTO.StoreID,
		StreetName: &bodyDTO.Data.StreetName,
		LocationID: &bodyDTO.Data.LocationID,
	}
	if bodyDTO.Data.LocationID == 0 {
		dataUpdate.LocationID = nil
	}

	address, err = instance.UpdateByID(paramDTO.ID, dataUpdate)
	if err != nil {
		instance.Debug(
			"UPDATE_BY_id.CreateOne",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.UnprocessableEntityException(err.Error()))
	}

	return address
}

func (instance REST) DELETE_BY_id(
	ctx gogo.Context,
	tokenClaimsDTO shared.TokenClaimsDTO,
	paramDTO dtos.DELETE_BY_id_Param,
) gogo.Map {
	_, err := instance.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		instance.Debug(
			"DELETE_BY_id.FindOneBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return gogo.Map{
			"deleted": false,
		}
	}

	if err = instance.DeleteByID(paramDTO.ID); err != nil {
		instance.Debug(
			"DELETE_BY_id.DeleteByID",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return gogo.Map{
			"deleted": false,
		}
	}

	return gogo.Map{
		"deleted": true,
	}
}
