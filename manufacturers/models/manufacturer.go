package models

import (
	"time"

	storeModels "github.com/dangduoc08/ecommerce-api/stores/models"
)

type Manufacturer struct {
	ID        uint              `json:"id" gorm:"primaryKey"`
	StoreID   uint              `json:"-"  gorm:"not null"`
	Store     storeModels.Store `json:"-" gorm:"foreignKey:StoreID"`
	Name      string            `json:"name" gorm:"not null"`
	Slug      string            `json:"slug" gorm:"unique;not null"`
	Logo      *string           `json:"logo"`
	CreatedAt time.Time         `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt time.Time         `json:"updated_at" gorm:"autoUpdateTime:true"`
}
