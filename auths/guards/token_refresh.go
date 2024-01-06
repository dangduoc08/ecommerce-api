package guards

import (
	"context"
	"fmt"
	"strings"

	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenRefresh struct {
	config.ConfigService
	JWTRefreshTokenAPIKey string
}

func (self TokenRefresh) NewGuard() TokenRefresh {
	self.JWTRefreshTokenAPIKey = self.Get("JWT_REFRESH_TOKEN_KEY").(string)

	return self
}

func (self TokenRefresh) CanActivate(c gooh.Context) bool {
	refreshTokenCookie, err := c.Cookie("refresh_token")
	if err != nil {
		return false
	}
	refreshToken := strings.Replace(refreshTokenCookie.Value, "Bearer ", "", 1)
	if refreshToken == "" {
		return false
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(self.JWTRefreshTokenAPIKey), nil
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
