package providers

import (
	permissionModels "github.com/dangduoc08/ecommerce-api/permissions/models"
)

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
	Permissions permissionModels.Permissions
}

type Update struct {
	Name        string
	StoreID     uint
	Permissions permissionModels.Permissions
}
