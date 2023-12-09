package users

import (
	"errors"
	"fmt"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/dangduoc08/ecommerce-api/db"
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

func (self UserProvider) signToken(claims jwt.MapClaims, key string, expIn int) (string, error) {
	claims["exp"] = time.Now().Unix() + int64(expIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func (self UserProvider) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (self UserProvider) CheckDuplicateUser(data []map[string]string) error {
	errMsgs := []string{}

	for _, kv := range data {
		for k, v := range kv {
			var isExist uint
			sql := fmt.Sprintf("SELECT id FROM %v WHERE %v = '%v';", self.GetModelName(), k, v)
			self.DBProvider.DB.Raw(sql).Scan(&isExist)
			if isExist != 0 {
				errMsgs = append(errMsgs, fmt.Sprintf("%v: '%v' is duplicated", k, v))
			}
		}
	}

	if len(errMsgs) > 0 {
		return errors.New(strings.Join(errMsgs, ", "))
	}

	return nil
}

func (self UserProvider) CreateOneUser(dto CREATE_Body_DTO) (*User, error) {
	hash, err := self.HashPassword(dto.Data.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:  dto.Data.Username,
		Email:     dto.Data.Email,
		Hash:      hash,
		FirstName: dto.Data.FirstName,
		LastName:  dto.Data.LastName,
	}

	resp := self.DBProvider.DB.Create(user)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return user, nil
}

func (self UserProvider) GetOneUserBy(id uint) (*User, error) {
	user := &User{
		ID: id,
	}

	resp := self.DBProvider.DB.First(user)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return user, nil
}
