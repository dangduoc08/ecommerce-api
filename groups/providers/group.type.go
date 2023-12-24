package providers

import (
	permissionModels "github.com/dangduoc08/ecommerce-api/permissions/models"
)

type GroupQuery struct {
	IDs     []uint
	ID      uint
	StoreID uint
	Sort    string
	Order   string
	Limit   int
	Offset  int
}

type GroupCreation struct {
	Name        string
	StoreID     uint
	Permissions permissionModels.Permission
}

type GroupUpdate struct {
	Name        string
	StoreID     uint
	Permissions permissionModels.Permission
}
