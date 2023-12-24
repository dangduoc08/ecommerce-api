package providers

import (
	"time"

	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	groupModels "github.com/dangduoc08/ecommerce-api/groups/models"
	groupProviders "github.com/dangduoc08/ecommerce-api/groups/providers"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DBProvider    dbProviders.DB
	ConfigService config.ConfigService
	Logger        common.Logger
	GroupDB       groupProviders.GroupDB
}

func (self AuthHandler) NewProvider() core.Provider {
	return self
}

func (self AuthHandler) HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), nil
}

func (self AuthHandler) CheckHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (self AuthHandler) SignToken(claims jwt.MapClaims, key string, expIn int) (string, error) {
	claims["exp"] = time.Now().Unix() + int64(expIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

func (self AuthHandler) GetUserPermissions(grs []*groupModels.Group) []string {
	permissions := []string{}

	for _, gr := range grs {
		permissions = append(permissions, gr.Permissions...)
	}

	return utils.ArrToUnique[string](permissions)
}
