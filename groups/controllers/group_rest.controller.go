package controllers

import (
	"github.com/dangduoc08/ecommerce-api/groups/models"
	"github.com/dangduoc08/ecommerce-api/groups/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type GroupREST struct {
	common.REST
	common.Guard
	GroupDB providers.GroupDB
	Logger  common.Logger
}

func (self GroupREST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("groups")

	self.
		BindGuard(
			shared.AuthGuard{},
		)

	return self
}

func (self GroupREST) READ(
	c gooh.Context,
	queryDTO models.READ_Query,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
) []*models.Group {
	groups, err := self.GroupDB.FindManyBy(&providers.GroupQuery{
		StoreID: accessTokenPayloadDTO.StoreID,
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
	})

	if err != nil {
		self.Logger.Debug(
			"GroupREST.READ.GroupDB.FindManyBy",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return []*models.Group{}
	}

	return groups
}

func (self GroupREST) CREATE(
	c gooh.Context,
	bodyDTO models.CREATE_Body,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
) *models.Group {
	group, err := self.GroupDB.CreateOne(&providers.GroupCreation{
		Name:        bodyDTO.Data.Name,
		Permissions: bodyDTO.Data.Permissions,
		StoreID:     accessTokenPayloadDTO.StoreID,
	})

	if err != nil {
		self.Logger.Debug(
			"GroupREST.CREATE.GroupDB.CreateOne",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return group
}

func (self GroupREST) UPDATE_BY_id(
	c gooh.Context,
	paramDTO models.UPDATE_BY_id_Param,
	bodyDTO models.UPDATE_BY_id_Body,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
) *models.Group {
	group, err := self.GroupDB.FindOneBy(&providers.GroupQuery{
		ID:      paramDTO.ID,
		StoreID: accessTokenPayloadDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	group, err = self.GroupDB.UpdateByID(paramDTO.ID, &providers.GroupUpdate{
		Name:        bodyDTO.Data.Name,
		Permissions: bodyDTO.Data.Permissions,
		StoreID:     accessTokenPayloadDTO.StoreID,
	})
	if err != nil {
		self.Logger.Debug(
			"GroupREST.UPDATE_BY_id.GroupDB.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return group
}

func (self GroupREST) DELETE_BY_id(
	c gooh.Context,
	paramDTO models.DELETE_BY_id_Param,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
) gooh.Map {
	_, err := self.GroupDB.FindOneBy(&providers.GroupQuery{
		ID:      paramDTO.ID,
		StoreID: accessTokenPayloadDTO.StoreID,
	})

	if err != nil {
		self.Logger.Debug(
			"GroupREST.DELETE_BY_id.GroupDB.FindOneBy",
			"error", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return gooh.Map{
			"deleted": false,
		}
	}

	if err := self.GroupDB.DeleteByID(paramDTO.ID); err != nil {
		self.Logger.Debug(
			"GroupREST.DELETE_BY_id.GroupDB.DeleteByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return gooh.Map{
			"deleted": false,
		}
	}

	return gooh.Map{
		"deleted": true,
	}
}
