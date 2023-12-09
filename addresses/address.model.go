package addresses

import (
	"github.com/dangduoc08/ecommerce-api/locations"
)

type Address struct {
	ID         uint                `json:"id" gorm:"primaryKey"`
	StreetName string              `json:"street_name"`
	LocationID *uint               `json:"-" gorm:"nullable"`
	Location   *locations.Location `json:"location"`
	StoreID    uint                `json:"-"`
}
