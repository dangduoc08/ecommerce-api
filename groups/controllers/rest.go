package controllers

import (
	"github.com/dangduoc08/ecommerce-api/groups/dtos"
	"github.com/dangduoc08/ecommerce-api/groups/models"
	"github.com/dangduoc08/ecommerce-api/groups/providers"
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
		Prefix("groups")

	self.
		BindGuard(shared.AuthGuard{})

	return self
}

func (self REST) READ(
	ctx gooh.Context,
	queryDTO dtos.READ_Query,
	tokenClaimsDTO shared.TokenClaimsDTO,
) []*models.Group {
	groups, err := self.FindManyBy(&providers.Query{
		StoreID: tokenClaimsDTO.StoreID,
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
	})

	if err != nil {
		self.Logger.Debug(
			"READ.GroupDB.FindManyBy",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return []*models.Group{}
	}

	return groups
}

func (self REST) CREATE(
	ctx gooh.Context,
	bodyDTO dtos.CREATE_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Group {
	group, err := self.CreateOne(&providers.Creation{
		Name:        bodyDTO.Data.Name,
		Permissions: bodyDTO.Data.Permissions,
		StoreID:     tokenClaimsDTO.StoreID,
	})

	if err != nil {
		self.Logger.Debug(
			"CREATE.GroupDB.CreateOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return group
}

func (self REST) UPDATE_BY_id(
	ctx gooh.Context,
	paramDTO dtos.UPDATE_BY_id_Param,
	bodyDTO dtos.UPDATE_BY_id_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Group {
	group, err := self.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	group, err = self.UpdateByID(paramDTO.ID, &providers.Update{
		Name:        bodyDTO.Data.Name,
		Permissions: bodyDTO.Data.Permissions,
		StoreID:     tokenClaimsDTO.StoreID,
	})
	if err != nil {
		self.Logger.Debug(
			"UPDATE_BY_id.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return group
}

func (self REST) DELETE_BY_id(
	ctx gooh.Context,
	paramDTO dtos.DELETE_BY_id_Param,
	tokenClaimsDTO shared.TokenClaimsDTO,
) gooh.Map {
	_, err := self.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		self.Logger.Debug(
			"DELETE_BY_id.FindOneBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return gooh.Map{
			"deleted": false,
		}
	}

	if err := self.DeleteByID(paramDTO.ID); err != nil {
		self.Logger.Debug(
			"DELETE_BY_id.DeleteByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return gooh.Map{
			"deleted": false,
		}
	}

	return gooh.Map{
		"deleted": true,
	}
}
