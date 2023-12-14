package users

import (
	"errors"
	"fmt"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/dangduoc08/ecommerce-api/db"
	"github.com/dangduoc08/ecommerce-api/groups"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserProvider struct {
	DBProvider    db.DBProvider
	ConfigService config.ConfigService
	Logger        common.Logger
}

func (self UserProvider) NewProvider() core.Provider {
	return self
}

func (self UserProvider) GetModelName() string {
	return path.Base(reflect.TypeOf(self).PkgPath())
}

func (self UserProvider) HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func (self UserProvider) CheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (self UserProvider) SignToken(claims jwt.MapClaims, key string, expIn int) (string, error) {
	claims["exp"] = time.Now().Unix() + int64(expIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func (self UserProvider) CheckDuplicateUser(data []map[string]string) error {
	errMsgs := []string{}

	for _, kv := range data {
		for k, v := range kv {
			var userRec User
			self.DBProvider.DB.Where(fmt.Sprintf("%v = ?", k), fmt.Sprintf("%v", v)).First(&userRec)
			if userRec.ID != 0 {
				errMsgs = append(errMsgs, fmt.Sprintf("%v: '%v' is duplicated", k, v))
			}
		}
	}

	if len(errMsgs) > 0 {
		return errors.New(strings.Join(errMsgs, ", "))
	}

	return nil
}

func (self UserProvider) FindByID(id uint) (*User, error) {
	user := &User{
		ID: id,
	}

	resp := self.DBProvider.DB.First(user)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return user, nil
}

func (self UserProvider) FindOneBy(query *UserQuery) (*User, error) {
	user := &User{}

	if query.Username != "" {
		user.Username = query.Username
	}

	if query.Email != "" {
		user.Email = query.Email
	}

	resp := self.DBProvider.DB.Where(user).Preload("Groups").First(user)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return user, nil
}

func (self UserProvider) CreateOne(data *UserCreation) (*User, error) {
	hash, err := self.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		StoreID:   data.StoreID,
		Username:  data.Username,
		Email:     data.Email,
		Hash:      hash,
		FirstName: data.FirstName,
		LastName:  data.LastName,
	}

	resp := self.DBProvider.DB.Create(user)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return user, nil
}

func (self UserProvider) getUserPermissions(grs []groups.Group) []string {
	permissions := []string{}

	for _, gr := range grs {
		permissions = append(permissions, gr.Permissions...)
	}

	return utils.ArrToUnique[string](permissions)
}
