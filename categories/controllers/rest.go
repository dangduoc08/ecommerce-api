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
		BindGuard(
			shared.AuthGuard{},
		)

	return self
}

func (self REST) READ(
	c gooh.Context,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
	queryDTO dtos.READ_Query,
) any {
	categories, err := self.FindManyBy(&providers.Query{
		StoreID: accessTokenPayloadDTO.StoreID,
		Status:  models.CategoryStatus(queryDTO.Status),
		Sort:    queryDTO.Sort,
		Order:   queryDTO.Order,
		Limit:   queryDTO.Limit,
		Offset:  queryDTO.Offset,
	})

	if err != nil {
		self.Logger.Debug(
			"REST.READ.FindManyBy",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)

		return []*models.Category{}
	}

	return categories
}

func (self REST) READ_BY_id(
	c gooh.Context,
	paramDTO dtos.READ_BY_id_Param,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
) *models.Category {
	category, err := self.FindByID(paramDTO.ID)

	if err != nil {
		self.Logger.Debug(
			"REST.READ_BY_id.FindByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		return nil
	}

	return category
}

func (self REST) CREATE(
	c gooh.Context,
	bodyDTO dtos.CREATE_Body,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
) *models.Category {
	category, err := self.CreateOne(&providers.Creation{
		Name:              bodyDTO.Data.Name,
		Description:       bodyDTO.Data.Description,
		StoreID:           accessTokenPayloadDTO.StoreID,
		MetaTitle:         bodyDTO.Data.MetaTitle,
		MetaDescription:   bodyDTO.Data.MetaDescription,
		Slug:              bodyDTO.Data.Slug,
		Status:            models.CategoryStatus(bodyDTO.Data.Status),
		ParentCategoryIDs: bodyDTO.Data.ParentCategoryIDs,
	})

	if err != nil {
		self.Debug(
			"REST.CREATE.CreateOne",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return category
}

func (self REST) UPDATE_BY_id(
	c gooh.Context,
	paramDTO dtos.UPDATE_BY_id_Param,
	bodyDTO dtos.UPDATE_BY_id_Body,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
) *models.Category {
	_, err := self.FindOneBy(&providers.Query{
		ID:      paramDTO.ID,
		StoreID: accessTokenPayloadDTO.StoreID,
	})

	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	category, err := self.UpdateByID(paramDTO.ID, &providers.Update{
		StoreID:           accessTokenPayloadDTO.StoreID,
		Name:              bodyDTO.Data.Name,
		Description:       bodyDTO.Data.Description,
		MetaTitle:         bodyDTO.Data.MetaTitle,
		MetaDescription:   bodyDTO.Data.MetaDescription,
		Slug:              bodyDTO.Data.Slug,
		Status:            models.CategoryStatus(bodyDTO.Data.Status),
		ParentCategoryIDs: bodyDTO.Data.ParentCategoryIDs,
	})

	if err != nil {
		self.Logger.Debug(
			"REST.UPDATE_BY_id.UpdateByID",
			"message", err.Error(),
			"X-Request-ID", c.GetID(),
		)
		panic(exception.InternalServerErrorException(err.Error()))
	}

	return category
}

func (self REST) DELETE_BY_id(
	c gooh.Context,
	// bodyDTO models.CREATE_Body,
	accessTokenPayloadDTO shared.AccessTokenPayloadDTO,
) *models.Category {
	// group, err := self.GroupDB.CreateOne(&providers.GroupCreation{
	// 	Name:        bodyDTO.Data.Name,
	// 	Permissions: bodyDTO.Data.Permissions,
	// 	StoreID:     accessTokenPayloadDTO.StoreID,
	// })

	// if err != nil {
	// 	self.Logger.Debug(
	// 		"GroupREST.CREATE.GroupDB.CreateOne",
	// 		"message", err.Error(),
	// 		"X-Request-ID", c.GetID(),
	// 	)
	// 	panic(exception.InternalServerErrorException(err.Error()))
	// }

	return nil
}
