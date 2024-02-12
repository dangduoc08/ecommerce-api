package providers

import (
	categoryModels "github.com/dangduoc08/ecommerce-api/categories/models"
	manufacturerModels "github.com/dangduoc08/ecommerce-api/manufacturers/models"
	"github.com/dangduoc08/ecommerce-api/products/models"
)

type Query struct {
	ID           uint
	StoreID      uint
	CategoryID   uint
	Status       models.ProductStatus
	Manufacturer uint
	CategoryIDs  []uint
	Sort         string
	Order        string
	Limit        int
	Offset       int
}

type Creation struct {
	StoreID         uint
	Name            string
	Description     *string
	MetaTitle       string
	MetaDescription *string
	Slug            string
	Quantity        int
	SKU             string
	Height          float64
	Width           float64
	Length          float64
	Weight          float64
	Categories      []*categoryModels.Category
	Manufacturer    *manufacturerModels.Manufacturer
	VariantIDs      []uint
	Status          models.ProductStatus
	Images          models.ProductImages
}
