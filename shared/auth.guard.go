package shared

import (
	"context"
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthGuard struct {
	ConfigService   config.ConfigService
	JWTAccessAPIKey string
}

func (self AuthGuard) NewGuard() AuthGuard {
	self.JWTAccessAPIKey = self.ConfigService.Get("JWT_ACCESS_API_KEY").(string)

	return self
}

func (self AuthGuard) checkPermission(accessTo string, permissions []any) bool {
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

func (self AuthGuard) CanActivate(c gooh.Context) bool {
	accessTokenCookie, err := c.Cookie("access_token")
	if err != nil {
		return false
	}
	accessToken := strings.Replace(accessTokenCookie.Value, "Bearer ", "", 1)
	if accessToken == "" {
		return false
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(self.JWTAccessAPIKey), nil
	})

	if err != nil {
		panic(exception.UnauthorizedException(err.Error()))
	}

	if token.Claims != nil {
		matchedRoute := c.Method + c.GetRoute()
		ctxWithValue := context.WithValue(c.Context(), "tokenClaims", token.Claims.(jwt.MapClaims))
		c.Request = c.WithContext(ctxWithValue)

		return self.checkPermission(matchedRoute, token.Claims.(jwt.MapClaims)["permissions"].([]any))
	}

	return false
}
