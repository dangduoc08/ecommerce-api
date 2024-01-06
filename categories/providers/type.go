package providers

import "github.com/dangduoc08/ecommerce-api/categories/models"

type Query struct {
	ID         uint
	StoreID    uint
	CategoryID uint
	Status     models.CategoryStatus
	Sort       string
	Order      string
	Limit      int
	Offset     int
}

type Creation struct {
	StoreID           uint
	Name              string
	Description       *string
	MetaTitle         string
	MetaDescription   *string
	Slug              string
	Status            models.CategoryStatus
	ParentCategoryIDs []uint
}

type Update struct {
	ID                uint
	StoreID           uint
	Name              string
	Description       *string
	MetaTitle         string
	MetaDescription   *string
	Slug              string
	Status            models.CategoryStatus
	ParentCategoryIDs []uint
}
