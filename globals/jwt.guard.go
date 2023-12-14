package globals

import (
	"context"
	"fmt"
	"strings"

	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/golang-jwt/jwt/v5"
)

type AccessAPIGuard struct {
	ConfigService   config.ConfigService
	JWTAccessAPIKey string
}

func (self AccessAPIGuard) NewGuard() AccessAPIGuard {
	self.JWTAccessAPIKey = self.ConfigService.Get("JWT_ACCESS_API_KEY").(string)

	return self
}

func (self AccessAPIGuard) CanActivate(c gooh.Context) bool {
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
		ctxWithValue := context.WithValue(c.Context(), "tokenClaims", token.Claims.(jwt.MapClaims))
		c.Request = c.WithContext(ctxWithValue)
		return true
	}

	return false
}
