package stores

import (
	"time"

	"github.com/dangduoc08/ecommerce-api/admins/addresses"
	"github.com/dangduoc08/ecommerce-api/admins/users"
)

type StoreModel struct {
	ID          uint                     `json:"id" gorm:"primaryKey"`
	Name        string                   `json:"name" gorm:"not null"`
	Description *string                  `json:"description"`
	Phone       *string                  `json:"phone"`
	Email       *string                  `json:"email"`
	Addresses   []addresses.AddressModel `json:"-" gorm:"foreignKey:StoreID;references:ID"`
	Users       []users.UserModel        `json:"-" gorm:"foreignKey:StoreID;references:ID"`
	CreatedAt   time.Time                `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time                `json:"updated_at" gorm:"autoUpdateTime:true"`
}

type Update struct {
	Name        string
	Description *string
	Phone       *string
	Email       *string
}
