package controllers

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/auths/dtos"
	"github.com/dangduoc08/ecommerce-api/auths/guards"
	"github.com/dangduoc08/ecommerce-api/auths/providers"
	"github.com/dangduoc08/ecommerce-api/constants"
	userModels "github.com/dangduoc08/ecommerce-api/users/models"
	userProviders "github.com/dangduoc08/ecommerce-api/users/providers"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/golang-jwt/jwt/v5"
)

type REST struct {
	common.REST
	common.Guard
	providers.Cipher
	userProviders.DBHandler
	config.ConfigService
	common.Logger

	JWTAccessAPIExpIn    int
	JWTAccessAPIKey      string
	JWTRefreshTokenExpIn int
	JWTRefreshTokenKey   string
}

func (self REST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("auths")

	self.
		BindGuard(
			guards.TokenRefresh{},
			self.CREATE_tokens,
		)

	self.JWTAccessAPIKey = self.Get("JWT_ACCESS_API_KEY").(string)
	self.JWTAccessAPIExpIn = self.Get("JWT_ACCESS_API_EXP_IN").(int)
	self.JWTRefreshTokenKey = self.Get("JWT_REFRESH_TOKEN_KEY").(string)
	self.JWTRefreshTokenExpIn = self.Get("JWT_REFRESH_TOKEN_EXP_IN").(int)

	return self
}

func (self REST) CREATE_sessions(
	ctx gooh.Context,
	bodyDTO dtos.CREATE_sessions_Body,
) gooh.Map {
	user, err := self.FindOneBy(&userProviders.Query{
		Username: bodyDTO.Data.Username,
	})

	if err != nil {
		self.Error(
			"CREATE_sessions.FindOneBy",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	// check user should be active
	if user.Status != userModels.UserStatus(constants.USER_ACTIVE) {
		panic(exception.UnauthorizedException(fmt.Sprintf("user'status is %v", user.Status)))
	}

	if !self.CheckHash(bodyDTO.Data.Password, user.Hash) {
		panic(exception.UnauthorizedException("password not match"))
	}

	permissions := self.GetUserPermissions(user.Groups)

	accessToken, err := self.SignToken(
		jwt.MapClaims{
			"id":          user.ID,
			"store_id":    user.StoreID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"permissions": permissions,
		},
		self.JWTAccessAPIKey,
		self.JWTAccessAPIExpIn,
	)
	if err != nil {
		self.Error(
			"UserProvider.SignToken",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	refreshToken, err := self.SignToken(
		jwt.MapClaims{
			"id":       user.ID,
			"store_id": user.StoreID,
		},
		self.JWTRefreshTokenKey,
		self.JWTRefreshTokenExpIn,
	)
	if err != nil {
		self.Error(
			"UserProvider.SignToken",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return gooh.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gooh.Map{
			"id":          user.ID,
			"store_id":    user.StoreID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"permissions": permissions,
		},
	}
}

func (self REST) CREATE_tokens(
	ctx gooh.Context,
) gooh.Map {
	tokenClaims := ctx.Request.Context().Value("tokenClaims").(jwt.MapClaims)
	userID := uint(tokenClaims["id"].(float64))

	user, err := self.FindByID(userID)
	if err != nil {
		self.Error(
			"CREATE_tokens.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	if user.Status != userModels.UserStatus(constants.USER_ACTIVE) {
		panic(exception.UnauthorizedException(fmt.Sprintf("user'status is %v", user.Status)))
	}

	permissions := self.GetUserPermissions(user.Groups)

	accessToken, err := self.SignToken(
		jwt.MapClaims{
			"id":          user.ID,
			"store_id":    user.StoreID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"permissions": permissions,
		},
		self.JWTAccessAPIKey,
		self.JWTAccessAPIExpIn,
	)

	refreshTokenCookie, _ := ctx.Cookie("refresh_token")
	refreshToken := strings.Replace(refreshTokenCookie.Value, "Bearer ", "", 1)

	return gooh.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gooh.Map{
			"id":          user.ID,
			"store_id":    user.StoreID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"permissions": permissions,
		},
	}
}
