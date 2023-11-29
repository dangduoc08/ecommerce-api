package user

import (
	"github.com/dangduoc08/ecommerce-api/user/dtos"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type Controller struct {
	common.Rest
	common.Guard
	Provider Provider
}

func (controller Controller) NewController() core.Controller {
	controller.Prefix("v1").Prefix("users")

	controller.
		BindGuard(
			CreateGuard{},
			controller.CREATE_create,
		).
		BindGuard(
			SigninGuard{},
			controller.CREATE_signin,
		)

	return controller
}

func (controller Controller) CREATE_signin(
	ctx gooh.Context,
	dto dtos.CREATE_signin_Body_DTO,
) any {
	controller.Provider.genToken(gooh.Map{
		"name": "asda",
	})
	return gooh.Map{}
}

func (controller Controller) CREATE_create(
	dto dtos.CREATE_create_Body_DTO,
) User {
	user, err := controller.Provider.CreateOneUser(dto)
	if err != nil {
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return *user
}
