package groups

import (
	"time"

	"github.com/dangduoc08/ecommerce-api/admins/permissions"
)

type GroupModel struct {
	ID          uint                        `json:"id" gorm:"primaryKey"`
	Name        string                      `json:"name"`
	StoreID     uint                        `json:"-"`
	Permissions permissions.PermissionModel `json:"permissions" gorm:"type:TEXT"`
	CreatedAt   time.Time                   `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time                   `json:"updated_at" gorm:"autoUpdateTime:true"`
}

type Query struct {
	IDs     []uint
	ID      uint
	StoreID uint
	Sort    string
	Order   string
	Limit   int
	Offset  int
}

type Creation struct {
	Name        string
	StoreID     uint
	Permissions permissions.PermissionModel
}

type Update struct {
	Name        string
	StoreID     uint
	Permissions permissions.PermissionModel
}
