package models

import (
	"time"

	groupModels "github.com/dangduoc08/ecommerce-api/groups/models"
)

type UserStatus string

type User struct {
	ID        uint                 `json:"id" gorm:"primaryKey"`
	Username  string               `json:"username" gorm:"unique;not null"`
	Hash      string               `json:"-" gorm:"not null"`
	Email     string               `json:"email" gorm:"unique;not null"`
	FirstName string               `json:"first_name" gorm:"not null"`
	LastName  string               `json:"last_name" gorm:"not null"`
	Status    UserStatus           `json:"status" gorm:"not null;type:user_status;default:inactive"`
	Groups    []*groupModels.Group `json:"-" gorm:"many2many:users_groups"`
	StoreID   uint                 `json:"-"`
	CreatedAt time.Time            `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt time.Time            `json:"updated_at" gorm:"autoUpdateTime:true"`
}
