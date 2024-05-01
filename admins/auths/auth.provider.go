package auths

import (
	"time"

	"github.com/dangduoc08/ecommerce-api/admins/groups"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/modules/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthProvider struct {
	config.ConfigService
	common.Logger
}

func (instance AuthProvider) NewProvider() core.Provider {
	return instance
}

func (instance AuthProvider) GetUserPermissions(groups []*groups.GroupModel) []string {
	permissions := []string{}

	for _, group := range groups {
		permissions = append(permissions, group.Permissions...)
	}

	return utils.ArrToUnique(permissions)
}

func (instance AuthProvider) HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), nil
}

func (instance AuthProvider) CheckHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (instance AuthProvider) SignToken(claims jwt.MapClaims, key string, expIn int) (string, error) {
	claims["exp"] = time.Now().Unix() + int64(expIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}
