package providers

import (
	"time"

	dbProviders "github.com/dangduoc08/ecommerce-api/dbs/providers"
	groupModels "github.com/dangduoc08/ecommerce-api/groups/models"
	groupProviders "github.com/dangduoc08/ecommerce-api/groups/providers"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/modules/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Cipher struct {
	dbProviders.DB
	config.ConfigService
	common.Logger
	groupProviders.DBHandler
}

func (instance Cipher) NewProvider() core.Provider {
	return instance
}

func (instance Cipher) HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), nil
}

func (instance Cipher) CheckHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (instance Cipher) SignToken(claims jwt.MapClaims, key string, expIn int) (string, error) {
	claims["exp"] = time.Now().Unix() + int64(expIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

func (instance Cipher) GetUserPermissions(grs []*groupModels.Group) []string {
	permissions := []string{}

	for _, gr := range grs {
		permissions = append(permissions, gr.Permissions...)
	}

	return utils.ArrToUnique[string](permissions)
}
