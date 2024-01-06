package models

import (
	"time"

	addressModels "github.com/dangduoc08/ecommerce-api/addresses/models"
	userModels "github.com/dangduoc08/ecommerce-api/users/models"
)

type Store struct {
	ID          uint                    `json:"id" gorm:"primaryKey"`
	Name        string                  `json:"name" gorm:"not null"`
	Description string                  `json:"description"`
	Phone       string                  `json:"phone"`
	Email       string                  `json:"email"`
	Addresses   []addressModels.Address `json:"-" gorm:"foreignKey:StoreID;references:ID"`
	Users       []userModels.User       `json:"-" gorm:"foreignKey:StoreID;references:ID"`
	CreatedAt   time.Time               `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time               `json:"updated_at" gorm:"autoUpdateTime:true"`
}
