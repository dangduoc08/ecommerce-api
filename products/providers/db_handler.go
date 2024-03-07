package providers

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/dbs/providers"
	"github.com/dangduoc08/ecommerce-api/products/models"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/modules/config"
)

type DBHandler struct {
	dbProviders.DB
	config.ConfigService
}

func (instance DBHandler) NewProvider() core.Provider {
	return instance
}

func (instance DBHandler) FindByID(id uint) (*models.Product, error) {
	productRec := &models.Product{
		ID: id,
	}

	if err := instance.
		First(productRec).
		Preload("Categories").
		Preload("Variants").
		Preload("Manufacturer").
		Error; err != nil {
		return nil, err
	}

	return productRec, nil
}

func (instance DBHandler) FindManyBy(query *Query) ([]*models.Product, error) {
	productRecs := []*models.Product{}
	productQueries := map[string]any{}
	tx := instance.DB.DB

	if query.StoreID != 0 {
		productQueries["store_id"] = query.StoreID
	}

	if query.Status != "" {
		productQueries["status"] = query.Status
	}

	if query.Sort != "" {
		if query.Order == "" {
			query.Order = constants.ASC
		}
		tx = tx.Order(fmt.Sprintf("%v %v", query.Sort, query.Order))
	}

	if err := tx.
		Limit(query.Limit).
		Offset(query.Offset).
		Where(productQueries).
		Find(&productRecs).
		Error; err != nil {
		return []*models.Product{}, err
	}

	return productRecs, nil
}

func (instance DBHandler) CreateOne(data *Creation) (*models.Product, error) {
	productImages := models.ProductImages{}

	for _, image := range data.Images {
		productImages = append(productImages, models.ProductImage{
			URL:   image.URL,
			Order: image.Order,
		})
	}

	productRec := &models.Product{
		Name:            data.Name,
		Description:     data.Description,
		StoreID:         data.StoreID,
		MetaTitle:       data.MetaTitle,
		MetaDescription: data.MetaDescription,
		Slug:            data.Slug,
		Quantity:        data.Quantity,
		SKU:             data.SKU,
		Height:          data.Height,
		Width:           data.Width,
		Length:          data.Length,
		Weight:          data.Weight,
		Categories:      data.Categories,
		Manufacturer:    data.Manufacturer,
		// VariantIDs:      data.VariantIDs,
		Status: data.Status,
		Images: productImages,
	}

	if err := instance.
		Create(&productRec).
		Error; err != nil {
		return nil, err
	}

	return productRec, nil
}
