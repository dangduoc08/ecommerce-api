package controllers

import (
	"github.com/dangduoc08/ecommerce-api/manufacturers/dtos"
	"github.com/dangduoc08/ecommerce-api/manufacturers/models"
	"github.com/dangduoc08/ecommerce-api/manufacturers/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type REST struct {
	common.REST
	common.Guard
	common.Logger
	providers.DBHandler
}

func (self REST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("manufacturers")

	self.
		BindGuard(shared.AuthGuard{})

	return self
}

func (self REST) READ(
	ctx gooh.Context,
	tokenClaimsDTO shared.TokenClaimsDTO,
	queryDTO dtos.READ_Query,
) []*models.Manufacturer {
	manufacturers, err := self.FindManyBy(&providers.Query{
		StoreID: tokenClaimsDTO.StoreID,
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
	})

	if err != nil {
		self.Logger.Debug(
			"READ.FindManyBy",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)

		return []*models.Manufacturer{}
	}

	return manufacturers
}

func (self REST) READ_BY_id(
	ctx gooh.Context,
	paramDTO dtos.READ_BY_id_Param,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Manufacturer {
	manufacturer, err := self.FindByID(paramDTO.ID)

	if err != nil {
		self.Logger.Debug(
			"READ_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return nil
	}

	return manufacturer
}

func (self REST) CREATE(
	ctx gooh.Context,
	bodyDTO dtos.CREATE_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Manufacturer {
	manufacturer, err := self.CreateOne(&providers.Creation{
		StoreID: tokenClaimsDTO.StoreID,
		Name:    bodyDTO.Data.Name,
		Logo:    &bodyDTO.Data.Logo,
		Slug:    bodyDTO.Data.Slug,
	})

	if err != nil {
		self.Debug(
			"CREATE.CreateOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return manufacturer
}

func (self REST) UPDATE_BY_id(
	ctx gooh.Context,
	paramDTO dtos.UPDATE_BY_id_Param,
	bodyDTO dtos.UPDATE_BY_id_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Manufacturer {
	_, err := self.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	manufacturer, err := self.UpdateByID(paramDTO.ID, &providers.Update{
		Name: bodyDTO.Data.Name,
		Logo: &bodyDTO.Data.Logo,
		Slug: bodyDTO.Data.Slug,
	})

	if err != nil {
		self.Logger.Debug(
			"UPDATE_BY_id.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return manufacturer
}
