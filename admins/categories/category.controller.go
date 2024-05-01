package categories

import (
	"github.com/dangduoc08/ecommerce-api/admins/categories/dtos"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type CategoryController struct {
	common.REST
	common.Guard
	common.Logger
	CategoryProvider
}

func (instance CategoryController) NewController() core.Controller {
	instance.
		BindGuard(sharedLayers.AuthGuard{})

	return instance
}

func (instance CategoryController) READ_VERSION_1(
	ctx gogo.Context,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	queryDTO dtos.READ_Query_DTO,
) any {
	categories, err := instance.FindManyBy(&Query{
		StoreID: tokenClaimsDTO.StoreID,
		Status:  CategoryStatus(queryDTO.Status),
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

		return []*CategoryModel{}
	}

	return categories
}

func (instance CategoryController) READ_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.READ_BY_id_Param_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *CategoryModel {
	category, err := instance.FindByID(paramDTO.ID)

	if err != nil {
		instance.Logger.Debug(
			"READ_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return nil
	}

	return category
}

func (instance CategoryController) CREATE_VERSION_1(
	ctx gogo.Context,
	bodyDTO dtos.CREATE_Body,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *CategoryModel {
	category, err := instance.CreateOne(&Creation{
		Name:              bodyDTO.Data.Name,
		Description:       &bodyDTO.Data.Description,
		StoreID:           tokenClaimsDTO.StoreID,
		MetaTitle:         bodyDTO.Data.MetaTitle,
		MetaDescription:   &bodyDTO.Data.MetaDescription,
		Slug:              bodyDTO.Data.Slug,
		Status:            CategoryStatus(bodyDTO.Data.Status),
		ParentCategoryIDs: bodyDTO.Data.ParentCategoryIDs,
	})

	if err != nil {
		instance.Debug(
			"CREATE.CreateOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return category
}

func (instance CategoryController) UPDATE_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.UPDATE_BY_id_Param_DTO,
	bodyDTO dtos.UPDATE_BY_id_Body_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *CategoryModel {
	_, err := instance.FindOneBy(&Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	category, err := instance.UpdateByID(paramDTO.ID, &Update{
		StoreID:           tokenClaimsDTO.StoreID,
		Name:              bodyDTO.Data.Name,
		Description:       &bodyDTO.Data.Description,
		MetaTitle:         bodyDTO.Data.MetaTitle,
		MetaDescription:   &bodyDTO.Data.MetaDescription,
		Slug:              bodyDTO.Data.Slug,
		Status:            CategoryStatus(bodyDTO.Data.Status),
		ParentCategoryIDs: bodyDTO.Data.ParentCategoryIDs,
	})

	if err != nil {
		instance.Logger.Debug(
			"UPDATE_BY_id.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return category
}
