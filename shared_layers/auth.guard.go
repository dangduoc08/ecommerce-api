package sharedLayers

import (
	"context"
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/exception"
	"github.com/dangduoc08/gogo/modules/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthGuard struct {
	config.ConfigService
	JWTAccessAPIKey string
}

func (instance AuthGuard) NewGuard() AuthGuard {
	instance.JWTAccessAPIKey = instance.ConfigService.Get("JWT_ACCESS_API_KEY").(string)

	return instance
}

func (instance AuthGuard) checkPermission(accessTo string, permissions []any) bool {
	if utils.ArrIncludes(permissions, "*") {
		return true
	}

	for _, permission := range permissions {
		if accessTo == permission.(string) {
			return true
		}
	}

	return false
}

func (instance AuthGuard) CanActivate(c gogo.Context) bool {
	accessTokenCookie, err := c.Cookie("access_token")
	if err != nil {
		return false
	}
	accessToken := strings.Replace(accessTokenCookie.Value, constants.TOKEN_TYPE+" ", "", 1)
	if accessToken == "" {
		return false
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(instance.JWTAccessAPIKey), nil
	})

	if err != nil {
		panic(exception.UnauthorizedException(utils.Reason(err.Error())))
	}

	if token.Claims != nil {
		matchedRoute := c.Method + c.GetRoute()
		ctxWithValue := context.WithValue(c.Context(), TokenClaimContextKey("tokenClaims"), token.Claims.(jwt.MapClaims))
		c.Request = c.WithContext(ctxWithValue)

		return instance.checkPermission(matchedRoute, token.Claims.(jwt.MapClaims)["permissions"].([]any))
	}

	return false
}
