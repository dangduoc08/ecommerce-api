package stores

import (
	"github.com/dangduoc08/ecommerce-api/admins/stores/dtos"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type StoreController struct {
	common.REST
	common.Guard
	common.Interceptor
	common.Logger
	StoreProvider
}

func (instance StoreController) NewController() core.Controller {
	instance.
		BindGuard(
			sharedLayers.AuthGuard{},
		)

	return instance
}

func (instance StoreController) UPDATE_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.UPDATE_BY_id_Param_DTO,
	bodyDTO dtos.UPDATE_BY_id_Body_DTO,
) *StoreModel {
	_, err := instance.FindByID(paramDTO.ID)
	if err != nil {
		instance.Debug(
			"UPDATE_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	store, err := instance.UpdateByID(paramDTO.ID, &Update{
		Name:        bodyDTO.Data.Name,
		Description: bodyDTO.Data.Description,
		Phone:       bodyDTO.Data.Phone,
		Email:       bodyDTO.Data.Email,
	})

	if err != nil {
		instance.Debug(
			"UPDATE_BY_id.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return store
}
