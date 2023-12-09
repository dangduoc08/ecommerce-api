package users

import (
	"github.com/dangduoc08/ecommerce-api/globals"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/golang-jwt/jwt/v5"
)

type UserController struct {
	common.Rest
	common.Guard
	JWTAccessAPIExpIn    int
	JWTAccessAPIKey      string
	JWTRefreshTokenExpIn int
	JWTRefreshTokenKey   string
	UserProvider         UserProvider
	ConfigService        config.ConfigService
}

func (self UserController) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("users")

	self.JWTAccessAPIKey = self.ConfigService.Get("JWT_ACCESS_API_KEY").(string)
	self.JWTAccessAPIExpIn = self.ConfigService.Get("JWT_ACCESS_API_EXP_IN").(int)
	self.JWTRefreshTokenKey = self.ConfigService.Get("JWT_REFRESH_TOKEN_KEY").(string)
	self.JWTRefreshTokenExpIn = self.ConfigService.Get("JWT_REFRESH_TOKEN_EXP_IN").(int)

	self.
		BindGuard(
			globals.AccessAPIGuard{},
			self.CREATE,
		).
		BindGuard(
			SessionsGuard{},
			self.CREATE_sessions,
		)

	return self
}

func (self UserController) CREATE(dto CREATE_Body_DTO) User {
	user, err := self.UserProvider.CreateOneUser(dto)
	if err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return *user
}

func (self UserController) CREATE_sessions(dto CREATE_sessions_Body_DTO) any {
	user, err := self.UserProvider.GetOneUserBy(dto.Data.ID)
	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	accessToken, err := self.UserProvider.signToken(jwt.MapClaims{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}, self.JWTAccessAPIKey, self.JWTAccessAPIExpIn)
	if err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	refreshToken, err := self.UserProvider.signToken(jwt.MapClaims{
		"id": user.ID,
	}, self.JWTRefreshTokenKey, self.JWTRefreshTokenExpIn)
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
