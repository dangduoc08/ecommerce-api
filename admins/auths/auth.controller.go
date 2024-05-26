package auths

import (
	"net/http"

	"github.com/dangduoc08/ecommerce-api/admins/auths/dtos"
	mailConfigurations "github.com/dangduoc08/ecommerce-api/admins/mail_configurations"
	"github.com/dangduoc08/ecommerce-api/admins/stores"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/ecommerce-api/utils"

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
	UserProvider  users.UserProvider
	StoreProvider stores.StoreProvider

	JWTAccessAPIKey           string
	JWTAccessAPIExpIn         int
	JWTRefreshTokenKey        string
	JWTRefreshTokenExpIn      int
	MailConfigurationProvider mailConfigurations.MailConfigurationProvider
}

func (instance AuthController) NewController() core.Controller {
	instance.JWTAccessAPIKey = instance.Get("JWT_ACCESS_API_KEY").(string)
	instance.JWTAccessAPIExpIn = instance.Get("JWT_ACCESS_API_EXP_IN").(int)
	instance.JWTRefreshTokenKey = instance.Get("JWT_REFRESH_TOKEN_KEY").(string)
	instance.JWTRefreshTokenExpIn = instance.Get("JWT_REFRESH_TOKEN_EXP_IN").(int)

	instance.
		BindGuard(
			AuthGuard{},
			instance.MODIFY_refresh_token_VERSION_1,
			instance.MODIFY_recover_VERSION_1,
		)

	instance.
		BindGuard(
			sharedLayers.AuthGuard{},
			instance.READ_me_VERSION_1,
		)

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
		panic(exception.NotFoundException(utils.Reason(err.Error())))
	}

	// check user should be active
	if user.Status != users.UserStatus(constants.USER_ACTIVE) {
		panic(exception.UnauthorizedException(utils.Reason("user's status is", string(user.Status))))
	}

	if !instance.CheckHash(bodyDTO.Data.Password, user.Hash) {
		panic(exception.UnauthorizedException(utils.Reason("record not found")))
	}

	permissions := instance.GetUserPermissions(user.Groups)

	accessToken, err := instance.SignToken(
		jwt.MapClaims{
			"user_id":     user.ID,
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
			"user_id":  user.ID,
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
	}
}

func (instance AuthController) MODIFY_refresh_token_VERSION_1(
	ctx gogo.Context,
) gogo.Map {
	tokenClaims := ctx.Request.Context().Value(sharedLayers.TokenClaimContextKey("tokenClaims")).(jwt.MapClaims)
	userID := uint(tokenClaims["user_id"].(float64))

	user, err := instance.UserProvider.FindByID(userID)
	if err != nil {
		instance.Error(
			"MODIFY_refresh_token_VERSION_1.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(utils.Reason(err.Error())))
	}

	if user.Status != users.UserStatus(constants.USER_ACTIVE) {
		panic(exception.UnauthorizedException(utils.Reason("user's status is", string(user.Status))))
	}

	permissions := instance.GetUserPermissions(user.Groups)

	accessToken, _ := instance.SignToken(
		jwt.MapClaims{
			"user_id":     user.ID,
			"store_id":    user.StoreID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"permissions": permissions,
		},
		instance.JWTAccessAPIKey,
		instance.JWTAccessAPIExpIn,
	)

	return gogo.Map{
		"access": gogo.Map{
			"type":  constants.TOKEN_TYPE,
			"exp":   instance.JWTAccessAPIExpIn,
			"token": accessToken,
		},
	}
}

func (instance AuthController) MODIFY_request_password_reset_VERSION_1(
	ctx gogo.Context,
	bodyDTO dtos.MODIFY_reset_password_Body_DTO,
	headerDTO dtos.MODIFY_reset_password_Header_DTO,
) gogo.Map {
	userRec, err := instance.UserProvider.FindOneBy(
		&users.Query{
			Username: bodyDTO.Data.UserIdentity,
		},
		&users.Query{
			Email: bodyDTO.Data.UserIdentity,
		},
	)

	if err != nil {
		instance.Error(
			"MODIFY_request_password_reset_VERSION_1.UserProvider.FindOneBy",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)

		return gogo.Map{
			"requested": true,
		}
	}

	if userRec != nil && userRec.Status != constants.USER_ACTIVE {
		instance.Error(
			"MODIFY_request_password_reset_VERSION_1.UserProvider.FindOneBy",
			"user_status", userRec.Status,
			"X-Request-ID", ctx.GetID(),
		)

		return gogo.Map{
			"requested": true,
		}
	}

	storeRec, err := instance.StoreProvider.FindByID(userRec.StoreID)
	if err != nil {
		instance.Error(
			"MODIFY_request_password_reset_VERSION_1.StoreProvider.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)

		return gogo.Map{
			"requested": true,
		}
	}

	mailConfigurationRec, err := instance.MailConfigurationProvider.FindOneBy(
		&mailConfigurations.Query{
			StoreID: userRec.StoreID,
		},
	)
	if err != nil {
		instance.Error(
			"MODIFY_request_password_reset_VERSION_1.MailConfigurationProvider.FindOneBy",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)

		return gogo.Map{
			"requested": true,
		}
	}

	go instance.SendResetPasswordEmail(
		headerDTO.Origin,
		mailConfigurationRec,
		storeRec,
		userRec,
	)

	return gogo.Map{
		"requested": true,
	}
}

func (instance AuthController) MODIFY_recover_VERSION_1(
	ctx gogo.Context,
	bodyDTO dtos.MODIFY_recover_Body_DTO,
) (int, gogo.Map) {
	tokenClaims := ctx.Request.Context().Value(sharedLayers.TokenClaimContextKey("tokenClaims")).(jwt.MapClaims)
	userID := uint(tokenClaims["user_id"].(float64))

	user, err := instance.UserProvider.FindByID(userID)
	if err != nil {
		instance.Error(
			"MODIFY_recover_VERSION_1.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(utils.Reason(err.Error())))
	}

	if user.Status != users.UserStatus(constants.USER_ACTIVE) {
		panic(exception.UnauthorizedException(utils.Reason("user's status is", string(user.Status))))
	}

	user.Hash, err = instance.HashPassword(bodyDTO.Data.Password)
	if err != nil {
		instance.Error(
			"MODIFY_recover_VERSION_1.HashPassword",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(utils.Reason(err.Error())))
	}

	instance.UserProvider.ModifyOne(user)

	return http.StatusPermanentRedirect, gogo.Map{
		"recovered": true,
	}
}

func (instance AuthController) READ_me_VERSION_1(ctx gogo.Context) *users.UserModel {
	tokenClaims := ctx.Request.Context().Value(sharedLayers.TokenClaimContextKey("tokenClaims")).(jwt.MapClaims)
	userID := uint(tokenClaims["user_id"].(float64))

	user, err := instance.UserProvider.FindByID(userID)
	if err != nil {
		instance.Error(
			"READ_me_VERSION_1.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(utils.Reason(err.Error())))
	}

	if user.Status != users.UserStatus(constants.USER_ACTIVE) {
		panic(exception.UnauthorizedException(utils.Reason("user's status is", string(user.Status))))
	}

	return user
}
