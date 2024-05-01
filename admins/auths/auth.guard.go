package auths

import (
	"context"
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/constants"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/exception"
	"github.com/dangduoc08/gogo/modules/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthGuard struct {
	config.ConfigService
	JWTRefreshTokenAPIKey string
}

func (instance AuthGuard) NewGuard() AuthGuard {
	instance.JWTRefreshTokenAPIKey = instance.Get("JWT_REFRESH_TOKEN_KEY").(string)

	return instance
}

func (instance AuthGuard) CanActivate(c gogo.Context) bool {
	refreshTokenCookie, err := c.Cookie("refresh_token")
	if err != nil {
		return false
	}
	refreshToken := strings.Replace(refreshTokenCookie.Value, constants.TOKEN_TYPE+" ", "", 1)
	if refreshToken == "" {
		return false
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(instance.JWTRefreshTokenAPIKey), nil
	})

	if err != nil {
		panic(exception.UnauthorizedException(err.Error()))
	}

	if token.Claims != nil {
		ctxWithValue := context.WithValue(c.Context(), sharedLayers.TokenClaimContextKey("tokenClaims"), token.Claims.(jwt.MapClaims))
		c.Request = c.WithContext(ctxWithValue)

		return true
	}

	return false
}
