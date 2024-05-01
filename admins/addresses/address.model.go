package addresses

import (
	"time"

	"github.com/dangduoc08/ecommerce-api/admins/locations"
)

type AddressModel struct {
	ID         uint                     `json:"id" gorm:"primaryKey"`
	StreetName *string                  `json:"street_name"`
	LocationID *uint                    `json:"-" gorm:"nullable"`
	Location   *locations.LocationModel `json:"location"`
	StoreID    uint                     `json:"-"`
	CreatedAt  time.Time                `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt  time.Time                `json:"updated_at" gorm:"autoUpdateTime:true"`
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
	StoreID    uint
	StreetName *string
	LocationID *uint
}

type Update struct {
	StoreID    uint
	StreetName *string
	LocationID *uint
}
