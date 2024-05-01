package users

import (
	"github.com/dangduoc08/ecommerce-api/admins/users/dtos"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
	"github.com/dangduoc08/gogo/modules/config"
)

type UserController struct {
	common.REST
	common.Guard
	common.Logger
	config.ConfigService
	UserProvider
}

func (instance UserController) NewController() core.Controller {
	instance.
		BindGuard(sharedLayers.AuthGuard{})

	return instance
}

func (instance UserController) READ_VERSION_1(
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	queryDTO dtos.READ_Query_DTO,
	ctx gogo.Context,
) []*UserModel {
	users, err := instance.FindManyBy(&Query{
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
		return []*UserModel{}
	}

	return users
}

func (instance UserController) CREATE_VERSION_1(
	bodyDTO dtos.CREATE_Body_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	ctx gogo.Context,
) *UserModel {
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

	dataCreation := &Creation{
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

func (instance UserController) MODIFY_statuses_OF_BY_id_VERSION_1(
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	paramDTO dtos.MODIFY_statuses_OF_BY_id_Param_DTO,
	bodyDTO dtos.MODIFY_statuses_OF_BY_id_DTO,
	ctx gogo.Context,
) *UserModel {
	user, err := instance.FindByID(paramDTO.ID)
	if err != nil {
		instance.Error(
			"MODIFY_statuses_OF_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}
	user.Status = UserStatus(bodyDTO.Data.Status)

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
