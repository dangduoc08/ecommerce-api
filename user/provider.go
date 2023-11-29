package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/database"
	"github.com/dangduoc08/ecommerce-api/user/dtos"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Provider struct {
	DatabaseProvider database.Provider
	ConfigService    config.ConfigService
	JWTKey           string
}

func (provider Provider) NewProvider() core.Provider {
	provider.DatabaseProvider.CreateType("user_status", []string{ACTIVE, INACTIVE, SUSPENDED})
	provider.DatabaseProvider.DB.AutoMigrate(&User{})
	provider.JWTKey = provider.ConfigService.Get("JWT_KEY").(string)

	return provider
}

func (provider Provider) hashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func (provider Provider) genToken(data any) string {
	var (
		// key *ecdsa.PrivateKey
		t *jwt.Token
		s string
	)

	t = jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"iss": "my-auth-server",
			"sub": "john",
			"foo": 2,
		})
	s, _ = t.SignedString(provider.JWTKey)

	return s
}

func (provider Provider) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (provider Provider) CheckDuplicateUser(data []map[string]string) error {
	errMsgs := []string{}

	for _, kv := range data {
		for k, v := range kv {
			var isExist uint
			sql := fmt.Sprintf("SELECT id FROM users WHERE %v = '%v';", k, v)
			provider.DatabaseProvider.DB.Raw(sql).Scan(&isExist)
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

func (provider Provider) CreateOneUser(dto dtos.CREATE_create_Body_DTO) (*User, error) {
	hash, err := provider.hashPassword(dto.Data.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:  dto.Data.Username,
		Email:     dto.Data.Email,
		Hash:      hash,
		FirstName: dto.Data.FirstName,
		LastName:  dto.Data.LastName,
		Status:    INACTIVE,
	}

	resp := provider.DatabaseProvider.DB.Create(user)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return user, nil
}
