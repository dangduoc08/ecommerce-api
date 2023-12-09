package stores

import (
	"time"

	"github.com/dangduoc08/ecommerce-api/addresses"
	"github.com/dangduoc08/ecommerce-api/users"
)

type Store struct {
	ID          uint                `json:"id" gorm:"primaryKey"`
	Name        string              `json:"name" gorm:"not null"`
	Description string              `json:"description"`
	Phone       string              `json:"phone"`
	Email       string              `json:"email"`
	Addresses   []addresses.Address `json:"-" gorm:"foreignKey:StoreID;references:ID"`
	Users       []users.User        `json:"-" gorm:"foreignKey:StoreID;references:ID"`
	CreatedAt   time.Time           `json:"created_at" gorm:"autoCreateTime:true;"`
	UpdatedAt   time.Time           `json:"updated_at" gorm:"autoUpdateTime:true"`
}
