package user

import (
	"github.com/dangduoc08/ecommerce-api/user/dtos"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/golang-jwt/jwt/v5"
)

type Controller struct {
	common.Rest
	common.Guard
	Provider             Provider
	ConfigService        config.ConfigService
	JWTAccessAPIExpIn    int
	JWTRefreshTokenExpIn int
}

func (controller Controller) NewController() core.Controller {
	controller.Prefix("v1").Prefix("users")
	controller.JWTAccessAPIExpIn = controller.ConfigService.Get("JWT_ACCESS_API_EXP_IN").(int)
	controller.JWTRefreshTokenExpIn = controller.ConfigService.Get("JWT_REFRESH_TOKEN_EXP_IN").(int)

	controller.
		BindGuard(
			CreateGuard{},
			controller.CREATE_create,
		).
		BindGuard(
			SigninGuard{},
			controller.CREATE_signin,
		)

	return controller
}

func (controller Controller) CREATE_signin(
	ctx gooh.Context,
	dto dtos.CREATE_signin_Body_DTO,
) any {
	user, err := controller.Provider.GetOneUserBy(dto.Data.ID)
	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	accessToken, err := controller.Provider.signToken(jwt.MapClaims{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}, controller.JWTAccessAPIExpIn)
	if err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	refreshToken, err := controller.Provider.signToken(jwt.MapClaims{
		"id": user.ID,
	}, controller.JWTRefreshTokenExpIn)
	if err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return gooh.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gooh.Map{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		},
	}
}

func (controller Controller) CREATE_create(
	dto dtos.CREATE_create_Body_DTO,
) User {
	user, err := controller.Provider.CreateOneUser(dto)
	if err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return *user
}
