package groups

import (
	"github.com/dangduoc08/ecommerce-api/admins/groups/dtos"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type GroupController struct {
	common.REST
	common.Guard
	common.Logger
	GroupProvider
}

func (instance GroupController) NewController() core.Controller {
	instance.
		BindGuard(sharedLayers.AuthGuard{})

	return instance
}

func (instance GroupController) READ_VERSION_1(
	ctx gogo.Context,
	queryDTO dtos.READ_Query_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) []*GroupModel {
	groups, err := instance.FindManyBy(&Query{
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
		return []*GroupModel{}
	}

	return groups
}

func (instance GroupController) CREATE_VERSION_1(
	ctx gogo.Context,
	bodyDTO dtos.CREATE_Body_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *GroupModel {
	group, err := instance.CreateOne(&Creation{
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

func (instance GroupController) UPDATE_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.UPDATE_BY_id_Param_DTO,
	bodyDTO dtos.UPDATE_BY_id_Body,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *GroupModel {
	_, err := instance.FindOneBy(&Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	group, err := instance.UpdateByID(paramDTO.ID, &Update{
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

func (instance GroupController) DELETE_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.DELETE_BY_id_Param_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) gogo.Map {
	_, err := instance.FindOneBy(&Query{
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
