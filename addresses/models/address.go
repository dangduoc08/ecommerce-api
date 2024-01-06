package models

import (
	"time"

	locationModels "github.com/dangduoc08/ecommerce-api/locations/models"
)

type Address struct {
	ID         uint                     `json:"id" gorm:"primaryKey"`
	StreetName *string                  `json:"street_name"`
	LocationID *uint                    `json:"-" gorm:"nullable"`
	Location   *locationModels.Location `json:"location"`
	StoreID    uint                     `json:"-"`
	CreatedAt  time.Time                `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt  time.Time                `json:"updated_at" gorm:"autoUpdateTime:true"`
}
