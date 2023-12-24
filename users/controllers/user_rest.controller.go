package controllers

import (
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/ecommerce-api/users/models"
	"github.com/dangduoc08/ecommerce-api/users/providers"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
)

type UserREST struct {
	common.REST
	common.Guard

	UserDB        providers.UserDB
	ConfigService config.ConfigService
	Logger        common.Logger
}

func (self UserREST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("users")

	self.BindGuard(shared.AuthGuard{})

	return self
}

func (self UserREST) READ(
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	queryDTO models.READ_Query,
	c gooh.Context,
) []*models.User {
	users, err := self.UserDB.FindManyBy(&providers.UserQuery{
		StoreID:  accessTokenPayloadDTO.StoreID,
		Statuses: queryDTO.Statuses,
		Sort:     queryDTO.Sort,
		Order:    queryDTO.Order,
		Limit:    queryDTO.Limit,
		Offset:   queryDTO.Offset,
	})

	if err != nil {
		self.Logger.Debug(
			"UserREST.READ.UserDB.FindManyBy",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return []*models.User{}
	}

	return users
}

func (self UserREST) CREATE(
	bodyDTO models.CREATE_Body,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	c gooh.Context,
) *models.User {
	dataCheckDuplication := []map[string]string{
		{
			"email": bodyDTO.Data.Email,
		},
		{
			"username": bodyDTO.Data.Username,
		},
	}

	if self.UserDB.IsDuplicated(dataCheckDuplication) {
		panic(exception.ConflictException("user's information has taken"))
	}

	dataCreation := &providers.UserCreation{
		StoreID:   accessTokenPayloadDTO.StoreID,
		Username:  bodyDTO.Data.Username,
		Password:  bodyDTO.Data.Password,
		Email:     bodyDTO.Data.Email,
		FirstName: bodyDTO.Data.FirstName,
		LastName:  bodyDTO.Data.LastName,
		GroupIDs:  bodyDTO.Data.GroupIDs,
	}
	user, err := self.UserDB.CreateOne(dataCreation)
	if err != nil {
		self.Logger.Error(
			"UserREST.CREATE.UserDB.CreateOne",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return user
}

func (self UserREST) MODIFY_BY_id(
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	bodyDTO models.MODIFY_BY_id_Body,
	c gooh.Context,
) *models.User {
	user, err := self.UserDB.FindByID(accessTokenPayloadDTO.ID)
	if err != nil {
		self.Logger.Error(
			"UserREST.MODIFY_BY_id.UserDB.FindByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}
	user.Status = models.UserStatus(bodyDTO.Data.Status)

	user, err = self.UserDB.ModifyOne(user)
	if err != nil {
		self.Logger.Error(
			"UserREST.MODIFY_BY_id.UserDB.ModifyOne",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return user
}
