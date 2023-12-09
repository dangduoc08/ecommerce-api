package users

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/globals"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/golang-jwt/jwt/v5"
)

type UserController struct {
	common.REST
	common.Guard
	JWTAccessAPIExpIn    int
	JWTAccessAPIKey      string
	JWTRefreshTokenExpIn int
	JWTRefreshTokenKey   string
	UserProvider         UserProvider
	ConfigService        config.ConfigService
	Logger               common.Logger
}

func (self UserController) NewController() core.Controller {
	self.Prefix("v1").Prefix("users")
	self.BindGuard(globals.AccessAPIGuard{}, self.CREATE)

	self.JWTAccessAPIKey = self.ConfigService.Get("JWT_ACCESS_API_KEY").(string)
	self.JWTAccessAPIExpIn = self.ConfigService.Get("JWT_ACCESS_API_EXP_IN").(int)
	self.JWTRefreshTokenKey = self.ConfigService.Get("JWT_REFRESH_TOKEN_KEY").(string)
	self.JWTRefreshTokenExpIn = self.ConfigService.Get("JWT_REFRESH_TOKEN_EXP_IN").(int)

	return self
}

func (self UserController) CREATE(accessTokenDTO globals.AccessTokenDTO, bodyDTO CREATE_Body_DTO) *User {
	dataDuplication := []map[string]string{
		{
			"email": bodyDTO.Data.Email,
		},
		{
			"username": bodyDTO.Data.Username,
		},
	}
	err := self.UserProvider.CheckDuplicateUser(dataDuplication)
	if err != nil {
		errMsgs := []string{
			"Field: user.email or user.username",
			"Error: duplicate data",
		}
		panic(exception.ConflictException(errMsgs))
	}

	dataCreation := &UserCreation{
		StoreID:   accessTokenDTO.StoreID,
		Username:  bodyDTO.Data.Username,
		Password:  bodyDTO.Data.Password,
		Email:     bodyDTO.Data.Email,
		FirstName: bodyDTO.Data.FirstName,
		LastName:  bodyDTO.Data.LastName,
	}
	user, err := self.UserProvider.CreateOne(dataCreation)
	if err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return user
}

func (self UserController) CREATE_sessions(ctx gooh.Context, bodyDTO CREATE_sessions_Body_DTO) any {
	user, err := self.UserProvider.FindOneBy(&UserQuery{
		Username: bodyDTO.Data.Username,
	})

	if err != nil {
		self.Logger.Debug(
			"Error While Query",
			"error", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	// check user should be active
	if user.Status != UserStatus(ACTIVE) {
		panic(exception.UnauthorizedException(fmt.Sprintf("Field: user.status, Error: %s", user.Status)))
	}

	if !self.UserProvider.CheckHash(bodyDTO.Data.Password, user.Hash) {
		panic(exception.UnauthorizedException("Field: user.password, Error: not match"))
	}

	accessToken, err := self.UserProvider.SignToken(
		jwt.MapClaims{
			"id":         user.ID,
			"store_id":   user.StoreID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		},
		self.JWTAccessAPIKey,
		self.JWTAccessAPIExpIn,
	)
	if err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	refreshToken, err := self.UserProvider.SignToken(
		jwt.MapClaims{
			"id": user.ID,
		},
		self.JWTRefreshTokenKey,
		self.JWTRefreshTokenExpIn,
	)
	if err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return gooh.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gooh.Map{
			"id":         user.ID,
			"store_id":   user.StoreID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		},
	}
}
