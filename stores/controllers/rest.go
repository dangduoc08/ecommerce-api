package controllers

import (
	addressInterceptors "github.com/dangduoc08/ecommerce-api/addresses/interceptors"
	addressModels "github.com/dangduoc08/ecommerce-api/addresses/models"
	addressProviders "github.com/dangduoc08/ecommerce-api/addresses/providers"
	categoryInterceptors "github.com/dangduoc08/ecommerce-api/categories/interceptors"
	categoryProviders "github.com/dangduoc08/ecommerce-api/categories/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/ecommerce-api/stores/dtos"
	"github.com/dangduoc08/ecommerce-api/stores/models"
	"github.com/dangduoc08/ecommerce-api/stores/providers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type REST struct {
	common.REST
	common.Guard
	common.Interceptor
	common.Logger
	providers.DBHandler
	CategoryDBHandler categoryProviders.DBHandler
	AddressDBHandler  addressProviders.DBHandler
}

func (instance REST) NewController() core.Controller {
	instance.
		Prefix("v1").
		Prefix("stores")

	instance.
		BindGuard(
			shared.AuthGuard{},
			instance.UPDATE_BY_id,
		)

	instance.
		BindInterceptor(
			categoryInterceptors.SubCategoryTransformation{},
			instance.READ_categories_OF_BY_id,
		).
		BindInterceptor(
			addressInterceptors.AddressLocationTransformation{},
			instance.READ_addresses_OF_BY_id,
		)

	return instance
}

func (instance REST) READ_BY_id(
	ctx gogo.Context,
	paramDTO dtos.READ_BY_id_Param,
) *models.Store {
	store, err := instance.FindByID(paramDTO.ID)
	if err != nil {
		instance.Debug(
			"READ_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	return store
}

func (instance REST) READ_categories_OF_BY_id(
	ctx gogo.Context,
	paramDTO dtos.READ_categories_OF_BY_id_Param,
	queryDTO dtos.READ_categories_OF_BY_id_Query,
) *[]map[string]any {
	menu, err := instance.CategoryDBHandler.FindManyAsMenu(&categoryProviders.Query{
		StoreID:    paramDTO.ID,
		CategoryID: queryDTO.CategoryID,
	})

	if err != nil {
		instance.Debug(
			"READ_categories_OF_BY_id.FindManyAsMenu",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return &[]map[string]any{}
	}

	return menu
}

func (instance REST) READ_addresses_OF_BY_id(
	ctx gogo.Context,
	paramDTO dtos.READ_addresses_OF_BY_id_Param,
	queryDTO dtos.READ_addresses_OF_BY_id_Query,
) []*addressModels.Address {
	addresses, err := instance.AddressDBHandler.FindManyBy(&addressProviders.Query{
		StoreID: paramDTO.ID,
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
	})

	if err != nil {
		instance.Debug(
			"READ_addresses_OF_BY_id.AddressDBHandler.FindManyBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return []*addressModels.Address{}
	}

	return addresses
}

func (instance REST) UPDATE_BY_id(
	ctx gogo.Context,
	paramDTO dtos.UPDATE_BY_id_Param,
	bodyDTO dtos.UPDATE_BY_id_Body,
) *models.Store {
	_, err := instance.FindByID(paramDTO.ID)
	if err != nil {
		instance.Debug(
			"UPDATE_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	store, err := instance.UpdateByID(paramDTO.ID, &providers.Update{
		Name:        bodyDTO.Data.Name,
		Description: &bodyDTO.Data.Description,
		Phone:       &bodyDTO.Data.Phone,
		Email:       &bodyDTO.Data.Email,
	})

	if err != nil {
		instance.Debug(
			"UPDATE_BY_id.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return store
}
