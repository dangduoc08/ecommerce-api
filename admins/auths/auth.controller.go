package auths

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/admins/auths/dtos"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"

	"github.com/dangduoc08/ecommerce-api/admins/users"
	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
	"github.com/dangduoc08/gogo/modules/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	common.REST
	common.Guard
	common.Logger
	config.ConfigService

	AuthProvider
	UserProvider users.UserProvider

	JWTAccessAPIExpIn    int
	JWTAccessAPIKey      string
	JWTRefreshTokenExpIn int
	JWTRefreshTokenKey   string
}

func (instance AuthController) NewController() core.Controller {
	instance.
		BindGuard(
			AuthGuard{},
			instance.CREATE_tokens_VERSION_1,
		)

	instance.JWTAccessAPIKey = instance.Get("JWT_ACCESS_API_KEY").(string)
	instance.JWTAccessAPIExpIn = instance.Get("JWT_ACCESS_API_EXP_IN").(int)
	instance.JWTRefreshTokenKey = instance.Get("JWT_REFRESH_TOKEN_KEY").(string)
	instance.JWTRefreshTokenExpIn = instance.Get("JWT_REFRESH_TOKEN_EXP_IN").(int)

	return instance
}

func (instance AuthController) CREATE_sessions_VERSION_1(
	ctx gogo.Context,
	bodyDTO dtos.CREATE_sessions_Body_DTO,
) gogo.Map {
	user, err := instance.UserProvider.FindOneBy(&users.Query{
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
	if user.Status != users.UserStatus(constants.USER_ACTIVE) {
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
		"access": gogo.Map{
			"type":  constants.TOKEN_TYPE,
			"exp":   instance.JWTAccessAPIExpIn,
			"token": accessToken,
		},
		"refresh": gogo.Map{
			"type":  constants.TOKEN_TYPE,
			"exp":   instance.JWTRefreshTokenExpIn,
			"token": refreshToken,
		},
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

func (instance AuthController) CREATE_tokens_VERSION_1(
	ctx gogo.Context,
) gogo.Map {
	tokenClaims := ctx.Request.Context().Value(sharedLayers.TokenClaimContextKey("tokenClaims")).(jwt.MapClaims)
	userID := uint(tokenClaims["id"].(float64))

	user, err := instance.UserProvider.FindByID(userID)
	if err != nil {
		instance.Error(
			"CREATE_tokens.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	if user.Status != users.UserStatus(constants.USER_ACTIVE) {
		panic(exception.UnauthorizedException(fmt.Sprintf("user'status is %v", user.Status)))
	}

	permissions := instance.GetUserPermissions(user.Groups)

	accessToken, _ := instance.SignToken(
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
	refreshToken := strings.Replace(refreshTokenCookie.Value, constants.TOKEN_TYPE+" ", "", 1)

	return gogo.Map{
		"access": gogo.Map{
			"type":  constants.TOKEN_TYPE,
			"exp":   instance.JWTAccessAPIExpIn,
			"token": accessToken,
		},
		"refresh": gogo.Map{
			"type":  constants.TOKEN_TYPE,
			"exp":   instance.JWTRefreshTokenExpIn,
			"token": refreshToken,
		},
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
