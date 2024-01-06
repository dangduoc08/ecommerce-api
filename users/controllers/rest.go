package controllers

import (
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/ecommerce-api/users/dtos"
	"github.com/dangduoc08/ecommerce-api/users/models"
	"github.com/dangduoc08/ecommerce-api/users/providers"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
)

type REST struct {
	common.REST
	common.Guard
	common.Logger
	providers.DBHandler
	providers.DBValidation
	config.ConfigService
}

func (self REST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("users")

	self.
		BindGuard(shared.AuthGuard{})

	return self
}

func (self REST) READ(
	tokenClaimsDTO shared.TokenClaimsDTO,
	queryDTO dtos.READ_Query,
	ctx gooh.Context,
) []*models.User {
	users, err := self.FindManyBy(&providers.Query{
		StoreID:  tokenClaimsDTO.StoreID,
		Statuses: queryDTO.Statuses,
		Sort:     queryDTO.Sort,
		Order:    queryDTO.Order,
		Limit:    queryDTO.Limit,
		Offset:   queryDTO.Offset,
	})

	if err != nil {
		self.Debug(
			"READ.FindManyBy",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return []*models.User{}
	}

	return users
}

func (self REST) CREATE(
	bodyDTO dtos.CREATE_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
	ctx gooh.Context,
) *models.User {
	dataCheckDuplication := []map[string]string{
		{
			"email": bodyDTO.Data.Email,
		},
		{
			"username": bodyDTO.Data.Username,
		},
	}

	if self.CheckDuplicated(dataCheckDuplication) {
		panic(exception.ConflictException("user's information has taken"))
	}

	dataCreation := &providers.Creation{
		StoreID:   tokenClaimsDTO.StoreID,
		Username:  bodyDTO.Data.Username,
		Password:  bodyDTO.Data.Password,
		Email:     bodyDTO.Data.Email,
		FirstName: bodyDTO.Data.FirstName,
		LastName:  bodyDTO.Data.LastName,
		GroupIDs:  bodyDTO.Data.GroupIDs,
	}
	user, err := self.CreateOne(dataCreation)
	if err != nil {
		self.Error(
			"CREATE.CreateOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return user
}

func (self REST) MODIFY_statuses_OF_BY_id(
	tokenClaimsDTO shared.TokenClaimsDTO,
	paramDTO dtos.MODIFY_statuses_OF_BY_id_Param,
	bodyDTO dtos.MODIFY_statuses_OF_BY_id,
	ctx gooh.Context,
) *models.User {
	user, err := self.FindByID(paramDTO.ID)
	if err != nil {
		self.Error(
			"MODIFY_statuses_OF_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}
	user.Status = models.UserStatus(bodyDTO.Data.Status)

	user, err = self.ModifyOne(user)
	if err != nil {
		self.Error(
			"MODIFY_statuses_OF_BY_id.ModifyOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return user
}
