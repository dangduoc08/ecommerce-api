package controllers

import (
	"github.com/dangduoc08/ecommerce-api/categories/dtos"
	"github.com/dangduoc08/ecommerce-api/categories/models"
	"github.com/dangduoc08/ecommerce-api/categories/providers"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
)

type REST struct {
	providers.DBHandler
	providers.DBValidation
	common.REST
	common.Guard
	common.Logger
}

func (self REST) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("categories")

	self.
		BindGuard(shared.AuthGuard{})

	return self
}

func (self REST) READ(
	ctx gooh.Context,
	tokenClaimsDTO shared.TokenClaimsDTO,
	queryDTO dtos.READ_Query,
) any {
	categories, err := self.FindManyBy(&providers.Query{
		StoreID: tokenClaimsDTO.StoreID,
		Status:  models.CategoryStatus(queryDTO.Status),
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
	})

	if err != nil {
		self.Logger.Debug(
			"READ.FindManyBy",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)

		return []*models.Category{}
	}

	return categories
}

func (self REST) READ_BY_id(
	ctx gooh.Context,
	paramDTO dtos.READ_BY_id_Param,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Category {
	category, err := self.FindByID(paramDTO.ID)

	if err != nil {
		self.Logger.Debug(
			"READ_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		return nil
	}

	return category
}

func (self REST) CREATE(
	ctx gooh.Context,
	bodyDTO dtos.CREATE_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Category {
	category, err := self.CreateOne(&providers.Creation{
		Name:              bodyDTO.Data.Name,
		Description:       &bodyDTO.Data.Description,
		StoreID:           tokenClaimsDTO.StoreID,
		MetaTitle:         bodyDTO.Data.MetaTitle,
		MetaDescription:   &bodyDTO.Data.MetaDescription,
		Slug:              bodyDTO.Data.Slug,
		Status:            models.CategoryStatus(bodyDTO.Data.Status),
		ParentCategoryIDs: bodyDTO.Data.ParentCategoryIDs,
	})

	if err != nil {
		self.Debug(
			"CREATE.CreateOne",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return category
}

func (self REST) UPDATE_BY_id(
	ctx gooh.Context,
	paramDTO dtos.UPDATE_BY_id_Param,
	bodyDTO dtos.UPDATE_BY_id_Body,
	tokenClaimsDTO shared.TokenClaimsDTO,
) *models.Category {
	_, err := self.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: tokenClaimsDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	category, err := self.UpdateByID(paramDTO.ID, &providers.Update{
		StoreID:           tokenClaimsDTO.StoreID,
		Name:              bodyDTO.Data.Name,
		Description:       &bodyDTO.Data.Description,
		MetaTitle:         bodyDTO.Data.MetaTitle,
		MetaDescription:   &bodyDTO.Data.MetaDescription,
		Slug:              bodyDTO.Data.Slug,
		Status:            models.CategoryStatus(bodyDTO.Data.Status),
		ParentCategoryIDs: bodyDTO.Data.ParentCategoryIDs,
	})

	if err != nil {
		self.Logger.Debug(
			"UPDATE_BY_id.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", ctx.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return category
}
