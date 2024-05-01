package products

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/dbs"

	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/modules/config"
)

type ProductProvider struct {
	dbs.DBProvider
	config.ConfigService
}

func (instance ProductProvider) NewProvider() core.Provider {
	return instance
}

func (instance ProductProvider) FindByID(id uint) (*ProductModel, error) {
	productRec := &ProductModel{
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

func (instance ProductProvider) FindManyBy(query *Query) ([]*ProductModel, error) {
	productRecs := []*ProductModel{}
	productQueries := map[string]any{}
	tx := instance.DB

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
		return []*ProductModel{}, err
	}

	return productRecs, nil
}

func (instance ProductProvider) CreateOne(data *Creation) (*ProductModel, error) {
	productImages := ProductImages{}

	for _, image := range data.Images {
		productImages = append(productImages, ProductImageModel{
			URL:   image.URL,
			Order: image.Order,
		})
	}

	productRec := &ProductModel{
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
