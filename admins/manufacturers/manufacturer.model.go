package manufacturers

import (
	"time"

	"github.com/dangduoc08/ecommerce-api/admins/stores"
)

type ManufacturerModel struct {
	ID        uint              `json:"id" gorm:"primaryKey"`
	StoreID   uint              `json:"-"  gorm:"not null"`
	Store     stores.StoreModel `json:"-" gorm:"foreignKey:StoreID"`
	Name      string            `json:"name" gorm:"not null"`
	Slug      string            `json:"slug" gorm:"unique;not null"`
	Logo      *string           `json:"logo"`
	CreatedAt time.Time         `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt time.Time         `json:"updated_at" gorm:"autoUpdateTime:true"`
}

type Query struct {
	ID      uint
	StoreID uint
	Sort    string
	Order   string
	Limit   int
	Offset  int
}

type Creation struct {
	StoreID uint
	Name    string
	Slug    string
	Logo    *string
}

type Update struct {
	Name string
	Slug string
	Logo *string
}
