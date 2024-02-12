package controllers

import (
	categoryProviders "github.com/dangduoc08/ecommerce-api/categories/providers"
	manufacturerProviders "github.com/dangduoc08/ecommerce-api/manufacturers/providers"
	"github.com/dangduoc08/ecommerce-api/products/dtos"
	"github.com/dangduoc08/ecommerce-api/products/models"
	"github.com/dangduoc08/ecommerce-api/products/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type REST struct {
	common.REST
	common.Guard
	common.Logger
	providers.DBHandler
	CategoryDBValidation  categoryProviders.DBValidation
	ManufacturerDBHandler manufacturerProviders.DBHandler
}

func (instance REST) NewController() core.Controller {
	instance.
		Prefix("v1").
		Prefix("products")

	instance.
		BindGuard(shared.AuthGuard{})

	return instance
}

func (instance REST) READ(
	ctx gooh.Context,
	tokenClaimsDTO shared.TokenClaimsDTO,
	queryDTO dtos.READ_Query,
) []*models.Product {
	products, err := instance.FindManyBy(&providers.Query{
		StoreID: tokenClaimsDTO.StoreID,
		Status:  models.ProductStatus(queryDTO.Status),
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

		return []*models.Product{}
	}

	return products
}

func (instance REST) READ_BY_id(
	ctx gooh.Context,
	paramDTO dtos.READ_BY_id_Param,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Product {
	product, err := instance.FindByID(paramDTO.ID)

	if err != nil {
		instance.Logger.Debug(
			"READ_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return nil
	}

	return product
}

func (instance REST) CREATE(
	ctx gooh.Context,
	bodyDTO dtos.CREATE_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Product {

	categories, err := instance.CategoryDBValidation.CheckParentCategories(bodyDTO.Data.CategoryIDs)
	if err != nil {
		instance.Debug(
			"CREATE.CategoryDBValidation.CheckParentCategories",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	manufacturer, err := instance.ManufacturerDBHandler.FindByID(bodyDTO.Data.ManufacturerID)
	if err != nil {
		instance.Debug(
			"CREATE.ManufacturerDBHandler.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	product, err := instance.CreateOne(&providers.Creation{
		Name:            bodyDTO.Data.Name,
		Description:     &bodyDTO.Data.Description,
		StoreID:         tokenClaimsDTO.StoreID,
		MetaTitle:       bodyDTO.Data.MetaTitle,
		MetaDescription: &bodyDTO.Data.MetaDescription,
		Slug:            bodyDTO.Data.Slug,
		Quantity:        bodyDTO.Data.Quantity,
		SKU:             bodyDTO.Data.SKU,
		Height:          bodyDTO.Data.Height,
		Width:           bodyDTO.Data.Width,
		Length:          bodyDTO.Data.Length,
		Weight:          bodyDTO.Data.Weight,
		Categories:      categories,
		Manufacturer:    manufacturer,
		VariantIDs:      bodyDTO.Data.VariantIDs,
		Status:          models.ProductStatus(bodyDTO.Data.Status),
		Images:          bodyDTO.Data.Images,
	})

	if err != nil {
		instance.Debug(
			"CREATE.CreateOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return product
}

func (instance REST) UPDATE_BY_id(
// ctx gooh.Context,
// paramDTO dtos.UPDATE_BY_id_Param,
// bodyDTO dtos.UPDATE_BY_id_Body,
// tokenClaimsDTO shared.TokenClaimsDTO,
) any {
	// _, err := instance.FindOneBy(&providers.Query{
	// 	ID:      paramDTO.ID,
	// 	StoreID: tokenClaimsDTO.StoreID,
	// })

	// if err != nil {
	// 	panic(exception.NotFoundException(err.Error()))
	// }

	// category, err := instance.UpdateByID(paramDTO.ID, &providers.Update{
	// 	StoreID:           tokenClaimsDTO.StoreID,
	// 	Name:              bodyDTO.Data.Name,
	// 	Description:       &bodyDTO.Data.Description,
	// 	MetaTitle:         bodyDTO.Data.MetaTitle,
	// 	MetaDescription:   &bodyDTO.Data.MetaDescription,
	// 	Slug:              bodyDTO.Data.Slug,
	// 	Status:            models.CategoryStatus(bodyDTO.Data.Status),
	// 	ParentCategoryIDs: bodyDTO.Data.ParentCategoryIDs,
	// })

	// if err != nil {
	// 	instance.Logger.Debug(
	// 		"UPDATE_BY_id.UpdateByID",
	// 		"message", err.Error(),
	// 		"X-Request-ID", ctx.GetID(),
	// 	)
	// 	panic(exception.InternalServerErrorException(err.Error()))
	// }

	// return category
	return nil
}
