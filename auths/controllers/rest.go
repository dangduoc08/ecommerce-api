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
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
	"github.com/dangduoc08/gogo/modules/config"
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

func (instance REST) NewController() core.Controller {
	instance.
		Prefix("v1").
		Prefix("auths")

	instance.
		BindGuard(
			guards.TokenRefresh{},
			instance.CREATE_tokens,
		)

	instance.JWTAccessAPIKey = instance.Get("JWT_ACCESS_API_KEY").(string)
	instance.JWTAccessAPIExpIn = instance.Get("JWT_ACCESS_API_EXP_IN").(int)
	instance.JWTRefreshTokenKey = instance.Get("JWT_REFRESH_TOKEN_KEY").(string)
	instance.JWTRefreshTokenExpIn = instance.Get("JWT_REFRESH_TOKEN_EXP_IN").(int)

	return instance
}

func (instance REST) CREATE_sessions(
	ctx gogo.Context,
	bodyDTO dtos.CREATE_sessions_Body,
) gogo.Map {
	user, err := instance.FindOneBy(&userProviders.Query{
		Username: bodyDTO.Data.Username,
	})

	if err != nil {
		instance.Error(
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

	if !instance.CheckHash(bodyDTO.Data.Password, user.Hash) {
		panic(exception.UnauthorizedException("password not match"))
	}

	permissions := instance.GetUserPermissions(user.Groups)

	accessToken, err := instance.SignToken(
		jwt.MapClaims{
			"id":          user.ID,
			"store_id":    user.StoreID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"permissions": permissions,
		},
		instance.JWTAccessAPIKey,
		instance.JWTAccessAPIExpIn,
	)
	if err != nil {
		instance.Error(
			"UserProvider.SignToken",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	refreshToken, err := instance.SignToken(
		jwt.MapClaims{
			"id":       user.ID,
			"store_id": user.StoreID,
		},
		instance.JWTRefreshTokenKey,
		instance.JWTRefreshTokenExpIn,
	)
	if err != nil {
		instance.Error(
			"UserProvider.SignToken",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return gogo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gogo.Map{
			"id":          user.ID,
			"store_id":    user.StoreID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"permissions": permissions,
		},
	}
}

func (instance REST) CREATE_tokens(
	ctx gogo.Context,
) gogo.Map {
	tokenClaims := ctx.Request.Context().Value("tokenClaims").(jwt.MapClaims)
	userID := uint(tokenClaims["id"].(float64))

	user, err := instance.FindByID(userID)
	if err != nil {
		instance.Error(
			"CREATE_tokens.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	if user.Status != userModels.UserStatus(constants.USER_ACTIVE) {
		panic(exception.UnauthorizedException(fmt.Sprintf("user'status is %v", user.Status)))
	}

	permissions := instance.GetUserPermissions(user.Groups)

	accessToken, err := instance.SignToken(
		jwt.MapClaims{
			"id":          user.ID,
			"store_id":    user.StoreID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"permissions": permissions,
		},
		instance.JWTAccessAPIKey,
		instance.JWTAccessAPIExpIn,
	)

	refreshTokenCookie, _ := ctx.Cookie("refresh_token")
	refreshToken := strings.Replace(refreshTokenCookie.Value, "Bearer ", "", 1)

	return gogo.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gogo.Map{
			"id":          user.ID,
			"store_id":    user.StoreID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"permissions": permissions,
		},
	}
}
