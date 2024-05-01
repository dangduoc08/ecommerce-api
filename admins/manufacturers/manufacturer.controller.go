package manufacturers

import (
	"github.com/dangduoc08/ecommerce-api/admins/manufacturers/dtos"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type ManufacturerController struct {
	common.REST
	common.Guard
	common.Logger
	ManufacturerProvider
}

func (instance ManufacturerController) NewController() core.Controller {
	instance.
		BindGuard(sharedLayers.AuthGuard{})

	return instance
}

func (instance ManufacturerController) READ_VERSION_1(
	ctx gogo.Context,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	queryDTO dtos.READ_Query_DTO,
) []*ManufacturerModel {
	manufacturers, err := instance.FindManyBy(&Query{
		StoreID: tokenClaimsDTO.StoreID,
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
	})

	if err != nil {
		instance.Logger.Debug(
			"READ.FindManyBy",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)

		return []*ManufacturerModel{}
	}

	return manufacturers
}

func (instance ManufacturerController) READ_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.READ_BY_id_Param_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *ManufacturerModel {
	manufacturer, err := instance.FindByID(paramDTO.ID)

	if err != nil {
		instance.Logger.Debug(
			"READ_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return nil
	}

	return manufacturer
}

func (instance ManufacturerController) CREATE_VERSION_1(
	ctx gogo.Context,
	bodyDTO dtos.CREATE_Body_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *ManufacturerModel {
	manufacturer, err := instance.CreateOne(&Creation{
		StoreID: tokenClaimsDTO.StoreID,
		Name:    bodyDTO.Data.Name,
		Logo:    &bodyDTO.Data.Logo,
		Slug:    bodyDTO.Data.Slug,
	})

	if err != nil {
		instance.Debug(
			"CREATE.CreateOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return manufacturer
}

func (instance ManufacturerController) UPDATE_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.UPDATE_BY_id_Param_DTO,
	bodyDTO dtos.UPDATE_BY_id_Body_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *ManufacturerModel {
	_, err := instance.FindOneBy(&Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	manufacturer, err := instance.UpdateByID(paramDTO.ID, &Update{
		Name: bodyDTO.Data.Name,
		Logo: &bodyDTO.Data.Logo,
		Slug: bodyDTO.Data.Slug,
	})

	if err != nil {
		instance.Logger.Debug(
			"UPDATE_BY_id.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return manufacturer
}
