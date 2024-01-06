package models

import (
	"time"

	permissionModels "github.com/dangduoc08/ecommerce-api/permissions/models"
)

type Group struct {
	ID          uint                        `json:"id" gorm:"primaryKey"`
	Name        string                      `json:"name"`
	StoreID     uint                        `json:"-"`
	Permissions permissionModels.Permission `json:"permissions" gorm:"type:TEXT"`
	CreatedAt   time.Time                   `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time                   `json:"updated_at" gorm:"autoUpdateTime:true"`
}
