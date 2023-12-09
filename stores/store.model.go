package stores

import (
	"time"

	"github.com/dangduoc08/ecommerce-api/addresses"
)

type Store struct {
	ID          uint                `json:"id" gorm:"primaryKey"`
	Name        string              `json:"name" gorm:"not null" json:"username"`
	Description string              `json:"description"`
	Phone       string              `json:"phone"`
	Email       string              `json:"email"`
	Addresses   []addresses.Address `json:"addresses" gorm:"many2many:store_addresses"`
	CreatedAt   time.Time           `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time           `json:"updated_at" gorm:"autoUpdateTime:true"`
}
