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
}

func (instance AuthGuard) NewGuard() AuthGuard {
	return instance
}

func (instance AuthGuard) CanActivate(c gogo.Context) bool {
	headerKey := ""
	tokenKey := ""

	switch c.URL.Path {
	case "/admins/auths/recover":
		headerKey = "recover_token"
		tokenKey = instance.Get("JWT_RECOVER_KEY").(string)
	case "/admins/auths/refresh_token":
		headerKey = "refresh_token"
		tokenKey = instance.Get("JWT_REFRESH_TOKEN_KEY").(string)
	default:
		return false
	}

	jwtToken := c.Header().Get(headerKey)
	if jwtToken == "" {
		return false
	}
	jwtToken = strings.Replace(jwtToken, constants.TOKEN_TYPE+" ", "", 1)
	if jwtToken == "" {
		return false
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenKey), nil
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
