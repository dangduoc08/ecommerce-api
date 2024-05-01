package seeds

import (
	"github.com/dangduoc08/ecommerce-api/admins/addresses"
	"github.com/dangduoc08/ecommerce-api/admins/auths"
	"github.com/dangduoc08/ecommerce-api/admins/categories"
	"github.com/dangduoc08/ecommerce-api/admins/groups"
	"github.com/dangduoc08/ecommerce-api/admins/locations"
	"github.com/dangduoc08/ecommerce-api/admins/manufacturers"
	"github.com/dangduoc08/ecommerce-api/admins/products"
	"github.com/dangduoc08/ecommerce-api/admins/stores"
	"github.com/dangduoc08/ecommerce-api/admins/users"
	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/modules/config"
)

type SeedProvider struct {
	dbs.DBProvider
	config.ConfigService
	common.Logger
	JWTAccessAPIKey    string
	JWTRefreshTokenKey string
	LocationProvider   locations.LocationProvider
	AuthProvider       auths.AuthProvider
}

func (instance SeedProvider) NewProvider() core.Provider {

	// Create user_status enum
	instance.CreateEnum(constants.USER_STATUS_FIELD_NAME, constants.USER_STATUSES)

	// Create category_status enum
	instance.CreateEnum(constants.CATEGORY_STATUS_FIELD_NAME, constants.CATEGORY_STATUSES)

	// Create category_status enum
	instance.CreateEnum(constants.PRODUCT_STATUS_FIELD_NAME, constants.PRODUCT_STATUSES)

	// Create tables
	err := instance.AutoMigrate(&locations.LocationModel{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&addresses.AddressModel{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&stores.StoreModel{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&groups.GroupModel{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&users.UserModel{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&categories.CategoryModel{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&manufacturers.ManufacturerModel{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&products.ProductModel{})
	if err != nil {
		panic(err)
	}

	// Seed locations
	totalLocations := instance.Count(constants.TABLE_LOCATION)
	if totalLocations == 0 {
		instance.LocationProvider.Seed(func(locationRec locations.LocationModel) {
			resp := instance.Create(&locationRec)
			if resp.Error != nil {
				instance.Logger.Debug("SeedLocations", "error", resp.Error)
			}
		})
	}

	// Seed users
	hash, err := instance.AuthProvider.HashPassword(instance.ConfigService.Get("PASSWORD").(string))
	if err != nil {
		panic(err)
	}
	user := users.UserModel{
		Username:  instance.ConfigService.Get("USERNAME").(string),
		Email:     instance.ConfigService.Get("EMAIL").(string),
		FirstName: instance.ConfigService.Get("FIRST_NAME").(string),
		LastName:  instance.ConfigService.Get("LAST_NAME").(string),
		Hash:      hash,
		Status:    users.UserStatus(constants.USER_ACTIVE),
		Groups: []*groups.GroupModel{
			{
				Name:        "Administrator",
				Permissions: []string{"*"},
				StoreID:     1,
			},
		},
	}

	// Seed stores
	totalStores := instance.Count(constants.TABLE_STORE)
	if totalStores == 0 {
		resp := instance.Create(&stores.StoreModel{
			Name:      "Demo",
			Addresses: []addresses.AddressModel{},
			Users:     []users.UserModel{user},
		})
		if resp.Error != nil {
			instance.Logger.Debug("SeedStores", "error", resp.Error)
		}
	}

	return instance
}
