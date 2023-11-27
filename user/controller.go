package user

import (
	"github.com/dangduoc08/ecommerce-api/database"
	"github.com/dangduoc08/ecommerce-api/user/dtos"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
)

type Controller struct {
	common.Rest
	Provider         Provider
	DatabaseProvider database.Provider
}

func (controller Controller) NewController() core.Controller {
	controller.Prefix("v1").Prefix("users")

	return controller
}

func (controller Controller) CREATE_signin(
	dto dtos.CREATE_signin_Body_DTO,
) core.Controller {
	return controller
}

func (controller Controller) CREATE(
	dto dtos.CREATE_Body_DTO,
) User {
	user := controller.Provider.CreateOneUser(dto)

	return *user
}
