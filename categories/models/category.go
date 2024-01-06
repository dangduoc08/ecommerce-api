package models

import (
	"time"

	storeModels "github.com/dangduoc08/ecommerce-api/stores/models"
)

type CategoryStatus string

type Category struct {
	ID               uint              `json:"id" gorm:"primaryKey"`
	StoreID          uint              `json:"-"  gorm:"not null"`
	Store            storeModels.Store `json:"-" gorm:"foreignKey:StoreID"`
	Name             string            `json:"name" gorm:"not null"`
	Description      string            `json:"description"`
	MetaTitle        string            `json:"meta_title" gorm:"not null"`
	MetaDescription  string            `json:"meta_description"`
	Slug             string            `json:"slug" gorm:"unique;not null"`
	Status           CategoryStatus    `json:"status" gorm:"not null;type:category_status;default:disabled"`
	ParentCategories []*Category       `json:"parent_categories" gorm:"many2many:categories_categories"`
	CreatedAt        time.Time         `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt        time.Time         `json:"updated_at" gorm:"autoUpdateTime:true"`
}
