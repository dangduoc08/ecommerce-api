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
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type REST struct {
	common.REST
	common.Guard
	common.Interceptor
	common.Logger
	DBHandler         providers.DBHandler
	CategoryDBHandler categoryProviders.DBHandler
	AddressDBHandler  addressProviders.DBHandler
}

func (self REST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("stores")

	self.BindGuard(
		shared.AuthGuard{},
		self.UPDATE_BY_id,
	)

	self.
		BindInterceptor(
			categoryInterceptors.ChildCategoryTransformation{},
			self.READ_categories_OF_BY_id,
		).
		BindInterceptor(
			addressInterceptors.AddressLocationTransformation{},
			self.READ_addresses_OF_BY_id,
		)

	return self
}

func (self REST) READ_BY_id(
	c gooh.Context,
	paramDTO dtos.READ_BY_id_Param,
) *models.Store {
	store, err := self.DBHandler.FindByID(paramDTO.ID)
	if err != nil {
		self.Debug(
			"READ_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	return store
}

func (self REST) READ_categories_OF_BY_id(
	c gooh.Context,
	paramDTO dtos.READ_categories_OF_BY_id_Param,
	queryDTO dtos.READ_categories_OF_BY_id_Query,
) *[]map[string]any {
	menu, err := self.CategoryDBHandler.FindManyAsMenu(&categoryProviders.Query{
		StoreID:    paramDTO.ID,
		CategoryID: queryDTO.CategoryID,
	})

	if err != nil {
		self.Debug(
			"READ_categories_OF_BY_id.FindManyAsMenu",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return &[]map[string]any{}
	}

	return menu
}

func (self REST) READ_addresses_OF_BY_id(
	c gooh.Context,
	paramDTO dtos.READ_addresses_OF_BY_id_Param,
	queryDTO dtos.READ_addresses_OF_BY_id_Query,
) []*addressModels.Address {
	addresses, err := self.AddressDBHandler.FindManyBy(&addressProviders.Query{
		StoreID: paramDTO.ID,
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
	})

	if err != nil {
		self.Debug(
			"READ_addresses_OF_BY_id.AddressDBHandler.FindManyBy",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return []*addressModels.Address{}
	}

	return addresses
}

func (self REST) UPDATE_BY_id(
	c gooh.Context,
	paramDTO dtos.UPDATE_BY_id_Param,
	bodyDTO dtos.UPDATE_BY_id_Body,
) *models.Store {
	store, err := self.DBHandler.FindByID(paramDTO.ID)
	if err != nil {
		self.Debug(
			"UPDATE_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	store, err = self.DBHandler.UpdateByID(paramDTO.ID, &providers.Update{
		Name:        bodyDTO.Data.Name,
		Description: bodyDTO.Data.Description,
		Phone:       bodyDTO.Data.Phone,
		Email:       bodyDTO.Data.Email,
	})

	if err != nil {
		self.Debug(
			"UPDATE_BY_id.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return store
}
