package controllers

import (
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/ecommerce-api/users/dtos"
	"github.com/dangduoc08/ecommerce-api/users/models"
	"github.com/dangduoc08/ecommerce-api/users/providers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
	"github.com/dangduoc08/gogo/modules/config"
)

type REST struct {
	common.REST
	common.Guard
	common.Logger
	providers.DBHandler
	providers.DBValidation
	config.ConfigService
}

func (instance REST) NewController() core.Controller {
	instance.
		Prefix("v1").
		Prefix("users")

	instance.
		BindGuard(shared.AuthGuard{})

	return instance
}

func (instance REST) READ(
	tokenClaimsDTO shared.TokenClaimsDTO,
	queryDTO dtos.READ_Query,
	ctx gogo.Context,
) []*models.User {
	users, err := instance.FindManyBy(&providers.Query{
		StoreID:  tokenClaimsDTO.StoreID,
		Statuses: queryDTO.Statuses,
		Sort:     queryDTO.Sort,
		Order:    queryDTO.Order,
		Limit:    queryDTO.Limit,
		Offset:   queryDTO.Offset,
	})

	if err != nil {
		instance.Debug(
			"READ.FindManyBy",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return []*models.User{}
	}

	return users
}

func (instance REST) CREATE(
	bodyDTO dtos.CREATE_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
	ctx gogo.Context,
) *models.User {
	dataCheckDuplication := []map[string]string{
		{
			"email": bodyDTO.Data.Email,
		},
		{
			"username": bodyDTO.Data.Username,
		},
	}

	if instance.CheckDuplicated(dataCheckDuplication) {
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
	user, err := instance.CreateOne(dataCreation)
	if err != nil {
		instance.Error(
			"CREATE.CreateOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return user
}

func (instance REST) MODIFY_statuses_OF_BY_id(
	tokenClaimsDTO shared.TokenClaimsDTO,
	paramDTO dtos.MODIFY_statuses_OF_BY_id_Param,
	bodyDTO dtos.MODIFY_statuses_OF_BY_id,
	ctx gogo.Context,
) *models.User {
	user, err := instance.FindByID(paramDTO.ID)
	if err != nil {
		instance.Error(
			"MODIFY_statuses_OF_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}
	user.Status = models.UserStatus(bodyDTO.Data.Status)

	user, err = instance.ModifyOne(user)
	if err != nil {
		instance.Error(
			"MODIFY_statuses_OF_BY_id.ModifyOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return user
}
