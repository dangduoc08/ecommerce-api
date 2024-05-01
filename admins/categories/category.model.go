package categories

import (
	"time"

	"github.com/dangduoc08/ecommerce-api/admins/stores"
)

type CategoryStatus string

type CategoryModel struct {
	ID               uint              `json:"id" gorm:"primaryKey"`
	StoreID          uint              `json:"-"  gorm:"not null"`
	Store            stores.StoreModel `json:"-" gorm:"foreignKey:StoreID"`
	Name             string            `json:"name" gorm:"not null"`
	Description      *string           `json:"description"`
	MetaTitle        string            `json:"meta_title" gorm:"not null"`
	MetaDescription  *string           `json:"meta_description"`
	Slug             string            `json:"slug" gorm:"unique;not null"`
	Status           CategoryStatus    `json:"status" gorm:"not null;type:category_status;default:disabled"`
	ParentCategories []*CategoryModel  `json:"parent_categories" gorm:"many2many:categories_categories"`
	CreatedAt        time.Time         `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt        time.Time         `json:"updated_at" gorm:"autoUpdateTime:true"`
}

type CategoryCategoryModel struct {
	CategoryID       uint
	ParentCategoryID uint
}

type Query struct {
	ID         uint
	StoreID    uint
	CategoryID uint
	Status     CategoryStatus
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
	Status            CategoryStatus
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
	Status            CategoryStatus
	ParentCategoryIDs []uint
}
