package products

import (
	"github.com/dangduoc08/ecommerce-api/admins/categories"
	"github.com/dangduoc08/ecommerce-api/admins/manufacturers"
	"github.com/dangduoc08/ecommerce-api/admins/products/dtos"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/exception"
)

type ProductController struct {
	common.REST
	common.Guard
	common.Logger
	ProductProvider
	CategoryProvider     categories.CategoryProvider
	ManufacturerProvider manufacturers.ManufacturerProvider
}

func (instance ProductController) NewController() core.Controller {
	instance.
		BindGuard(sharedLayers.AuthGuard{})

	return instance
}

func (instance ProductController) READ_VERSION_1(
	ctx gogo.Context,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
	queryDTO dtos.READ_Query_DTO,
) []*ProductModel {
	products, err := instance.FindManyBy(&Query{
		StoreID: tokenClaimsDTO.StoreID,
		Status:  ProductStatus(queryDTO.Status),
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

		return []*ProductModel{}
	}

	return products
}

func (instance ProductController) READ_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.READ_BY_id_Param_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *ProductModel {
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

func (instance ProductController) CREATE_VERSION_1(
	ctx gogo.Context,
	bodyDTO dtos.CREATE_Body_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) *ProductModel {

	categories, err := instance.CategoryProvider.CheckParentCategories(bodyDTO.Data.CategoryIDs)
	if err != nil {
		instance.Debug(
			"CREATE.CategoryDBValidation.CheckParentCategories",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	manufacturer, err := instance.ManufacturerProvider.FindByID(bodyDTO.Data.ManufacturerID)
	if err != nil {
		instance.Debug(
			"CREATE.ManufacturerDBHandler.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.NotFoundException(err.Error()))
	}

	product, err := instance.CreateOne(&Creation{
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
		Status:          ProductStatus(bodyDTO.Data.Status),
		Images: utils.ArrMap(bodyDTO.Data.Images, func(el dtos.CREATE_Body_Data_Image_DTO, idx int) ProductImageModel {
			return ProductImageModel{
				URL:   el.URL,
				Order: el.Order,
			}
		}),
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

func (instance ProductController) UPDATE_BY_id_VERSION_1(
	ctx gogo.Context,
	paramDTO dtos.UPDATE_BY_id_Param_DTO,
	bodyDTO dtos.UPDATE_BY_id_Body_DTO,
	tokenClaimsDTO sharedLayers.TokenClaimsDTO,
) any {
	// _, err := instance.FindOneBy(&Query{
	// 	ID:      paramDTO.ID,
	// 	StoreID: tokenClaimsDTO.StoreID,
	// })

	// if err != nil {
	// 	panic(exception.NotFoundException(err.Error()))
	// }

	// category, err := instance.UpdateByID(paramDTO.ID, &Update{
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
