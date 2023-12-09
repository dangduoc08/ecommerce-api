package globals

import (
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

func (accessAPIGuard AccessAPIGuard) NewGuard() AccessAPIGuard {
	accessAPIGuard.JWTAccessAPIKey = accessAPIGuard.ConfigService.Get("JWT_ACCESS_API_KEY").(string)

	return accessAPIGuard
}

func (accessAPIGuard AccessAPIGuard) CanActivate(ctx gooh.Context) bool {
	accessTokenCookie, err := ctx.Cookie("access_token")
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
		return []byte(accessAPIGuard.JWTAccessAPIKey), nil
	})

	if err != nil {
		panic(exception.UnauthorizedException(err.Error()))
	}

	if token.Claims != nil {
		return true
	}

	return false
}
