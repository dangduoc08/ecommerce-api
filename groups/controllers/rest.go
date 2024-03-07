package controllers

import (
	"github.com/dangduoc08/ecommerce-api/groups/dtos"
	"github.com/dangduoc08/ecommerce-api/groups/models"
	"github.com/dangduoc08/ecommerce-api/groups/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type REST struct {
	common.REST
	common.Guard
	common.Logger
	providers.DBHandler
}

func (instance REST) NewController() core.Controller {
	instance.
		Prefix("v1").
		Prefix("groups")

	instance.
		BindGuard(shared.AuthGuard{})

	return instance
}

func (instance REST) READ(
	ctx gogo.Context,
	queryDTO dtos.READ_Query,
	tokenClaimsDTO shared.TokenClaimsDTO,
) []*models.Group {
	groups, err := instance.FindManyBy(&providers.Query{
		StoreID: tokenClaimsDTO.StoreID,
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
	})

	if err != nil {
		instance.Logger.Debug(
			"READ.GroupDB.FindManyBy",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return []*models.Group{}
	}

	return groups
}

func (instance REST) CREATE(
	ctx gogo.Context,
	bodyDTO dtos.CREATE_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Group {
	group, err := instance.CreateOne(&providers.Creation{
		Name:        bodyDTO.Data.Name,
		Permissions: bodyDTO.Data.Permissions,
		StoreID:     tokenClaimsDTO.StoreID,
	})

	if err != nil {
		instance.Logger.Debug(
			"CREATE.GroupDB.CreateOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return group
}

func (instance REST) UPDATE_BY_id(
	ctx gogo.Context,
	paramDTO dtos.UPDATE_BY_id_Param,
	bodyDTO dtos.UPDATE_BY_id_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Group {
	group, err := instance.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	group, err = instance.UpdateByID(paramDTO.ID, &providers.Update{
		Name:        bodyDTO.Data.Name,
		Permissions: bodyDTO.Data.Permissions,
		StoreID:     tokenClaimsDTO.StoreID,
	})
	if err != nil {
		instance.Logger.Debug(
			"UPDATE_BY_id.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return group
}

func (instance REST) DELETE_BY_id(
	ctx gogo.Context,
	paramDTO dtos.DELETE_BY_id_Param,
	tokenClaimsDTO shared.TokenClaimsDTO,
) gogo.Map {
	_, err := instance.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		instance.Logger.Debug(
			"DELETE_BY_id.FindOneBy",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return gogo.Map{
			"deleted": false,
		}
	}

	if err := instance.DeleteByID(paramDTO.ID); err != nil {
		instance.Logger.Debug(
			"DELETE_BY_id.DeleteByID",
			"message", err.Error(),
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
