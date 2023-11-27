package user

import (
	"github.com/dangduoc08/ecommerce-api/database"
	"github.com/dangduoc08/ecommerce-api/user/dtos"
	"github.com/dangduoc08/gooh/core"
)

type Provider struct {
	DatabaseProvider database.Provider
}

func (provider Provider) NewProvider() core.Provider {
	provider.DatabaseProvider.DB.AutoMigrate(&User{})

	return provider
}

func (provider Provider) CreateOneUser(dto dtos.CREATE_Body_DTO) *User {
	userModel := &User{
		Username:  dto.Data.Username,
		Password:  dto.Data.Password,
		Email:     dto.Data.Email,
		Salt:      "12311",
		Firstname: dto.Data.Firstname,
		Lastname:  dto.Data.Lastname,
	}

	provider.DatabaseProvider.DB.Create(userModel)

	return userModel
}
