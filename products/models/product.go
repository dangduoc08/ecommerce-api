package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	categoryModels "github.com/dangduoc08/ecommerce-api/categories/models"
	manufacturerModels "github.com/dangduoc08/ecommerce-api/manufacturers/models"
	storeModels "github.com/dangduoc08/ecommerce-api/stores/models"
)

type ProductImage struct {
	URL   string `json:"url" bind:"url" validate:"http_url"`
	Order int    `json:"order" bind:"order" validate:"gte=0"`
}

type ProductImages []ProductImage

type ProductStatus string

func (productImages *ProductImages) Scan(src any) error {
	*productImages = ProductImages{}
	if srcStr, ok := src.(string); ok {
		return json.Unmarshal([]byte(srcStr), &productImages)
	}
	return nil
}

func (productImages ProductImages) Value() (driver.Value, error) {
	if len(productImages) == 0 {
		return nil, nil
	}

	jsonStrProductImages, err := json.Marshal(productImages)
	if err != nil {
		return nil, err
	}

	return string(jsonStrProductImages), nil
}

type Product struct {
	ID              uint                             `json:"id" gorm:"primaryKey"`
	StoreID         uint                             `json:"-"  gorm:"not null"`
	Store           storeModels.Store                `json:"-" gorm:"foreignKey:StoreID"`
	Name            string                           `json:"name" gorm:"not null"`
	Description     *string                          `json:"description"`
	MetaTitle       string                           `json:"meta_title" gorm:"not null"`
	MetaDescription *string                          `json:"meta_description"`
	Slug            string                           `json:"slug" gorm:"unique;not null"`
	Quantity        int                              `json:"quantity" gorm:"default:0"`
	SKU             string                           `json:"sku"`
	Height          float64                          `json:"height"  gorm:"default:0"`
	Width           float64                          `json:"width"  gorm:"default:0"`
	Length          float64                          `json:"length"  gorm:"default:0"`
	Weight          float64                          `json:"weight"  gorm:"default:0"`
	Categories      []*categoryModels.Category       `json:"categories" gorm:"many2many:products_categories"`
	Variants        []*Product                       `json:"variants" gorm:"many2many:products_products"`
	ManufacturerID  *uint                            `json:"-" gorm:"nullable"`
	Manufacturer    *manufacturerModels.Manufacturer `json:"manufacturer"`
	Status          ProductStatus                    `json:"status" gorm:"not null;type:product_status;default:disabled"`
	Images          ProductImages                    `json:"images" gorm:"type:TEXT"`
	CreatedAt       time.Time                        `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt       time.Time                        `json:"updated_at" gorm:"autoUpdateTime:true"`
}
