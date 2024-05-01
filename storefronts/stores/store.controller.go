package stores

import (
	"github.com/dangduoc08/ecommerce-api/admins/addresses"
	"github.com/dangduoc08/ecommerce-api/admins/categories"
	"github.com/dangduoc08/ecommerce-api/admins/stores"
	"github.com/dangduoc08/ecommerce-api/storefronts/stores/dtos"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type StoreController struct {
	common.REST
	common.Guard
	common.Interceptor
	common.Logger
	StoreProvider    stores.StoreProvider
	AddressProvider  addresses.AddressProvider
	CategoryProvider categories.CategoryProvider
}

func (instance StoreController) NewController() core.Controller {
	instance.
		BindInterceptor(
			categories.CategoryInterceptor{},
			instance.READ_categories_OF_BY_id_VERSION_1,
		).
		BindInterceptor(
			addresses.AddressInterceptor{},
			instance.READ_addresses_OF_BY_id_VERSION_1,
		)

	return instance
}

func (instance StoreController) READ_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.READ_BY_id_Param_DTO,
) *stores.StoreModel {
	store, err := instance.StoreProvider.FindByID(paramDTO.ID)
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

func (instance StoreController) READ_categories_OF_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.READ_categories_OF_BY_id_Param_DTO,
	queryDTO dtos.READ_categories_OF_BY_id_Query_DTO,
) *[]map[string]any {
	menu, err := instance.CategoryProvider.FindManyAsMenu(&categories.Query{
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

func (instance StoreController) READ_addresses_OF_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.READ_addresses_OF_BY_id_Param_DTO,
	queryDTO dtos.READ_addresses_OF_BY_id_Query_DTO,
) []*addresses.AddressModel {
	addressRecs, err := instance.AddressProvider.FindManyBy(&addresses.Query{
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
		return []*addresses.AddressModel{}
	}

	return addressRecs
}
