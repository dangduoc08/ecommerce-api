package addresses

import (
	"github.com/dangduoc08/ecommerce-api/admins/addresses/dtos"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type AddressController struct {
	common.Guard
	common.REST
	common.Interceptor
	common.Logger
	AddressProvider
}

func (instance AddressController) NewController() core.Controller {
	instance.
		BindGuard(sharedLayers.AuthGuard{})

	instance.
		BindInterceptor(AddressInterceptor{})

	return instance
}

func (instance AddressController) READ_VERSION_1(
	ctx gogo.Context,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	queryDTO dtos.READ_Query_DTO,
) []*AddressModel {
	addresses, err := instance.FindManyBy(&Query{
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
		ID:      queryDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		instance.Debug(
			"READ_BY_id.FindManyBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return nil
	}

	return addresses
}

func (instance AddressController) READ_BY_id_VERSION_1(
	ctx gogo.Context,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	paramDTO dtos.READ_BY_id_Param_DTO,
) *AddressModel {
	address, err := instance.FindOneBy(&Query{
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

func (instance AddressController) CREATE_VERSION_1(
	ctx gogo.Context,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	dto dtos.CREATE_Body_DTO,
) *AddressModel {
	dataCreation := &Creation{
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

func (instance AddressController) UPDATE_BY_id_VERSION_1(
	ctx gogo.Context,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	paramDTO dtos.UPDATE_BY_id_Param_DTO,
	bodyDTO dtos.UPDATE_BY_id_Body_DTO,
) *AddressModel {
	_, err := instance.FindOneBy(&Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	dataUpdate := &Update{
		StoreID:    tokenClaimsDTO.StoreID,
		StreetName: &bodyDTO.Data.StreetName,
		LocationID: &bodyDTO.Data.LocationID,
	}
	if bodyDTO.Data.LocationID == 0 {
		dataUpdate.LocationID = nil
	}

	address, err := instance.UpdateByID(paramDTO.ID, dataUpdate)
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

func (instance AddressController) DELETE_BY_id_VERSION_1(
	ctx gogo.Context,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	paramDTO dtos.DELETE_BY_id_Param_DTO,
) gogo.Map {
	_, err := instance.FindOneBy(&Query{
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
