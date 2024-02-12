package providers

import (
	addressModels "github.com/dangduoc08/ecommerce-api/addresses/models"
	authProviders "github.com/dangduoc08/ecommerce-api/auths/providers"
	categoryModels "github.com/dangduoc08/ecommerce-api/categories/models"
	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/dbs/providers"
	groupModels "github.com/dangduoc08/ecommerce-api/groups/models"
	locationModels "github.com/dangduoc08/ecommerce-api/locations/models"
	locationProviders "github.com/dangduoc08/ecommerce-api/locations/providers"
	manufacturerModels "github.com/dangduoc08/ecommerce-api/manufacturers/models"
	productModels "github.com/dangduoc08/ecommerce-api/products/models"
	storeModels "github.com/dangduoc08/ecommerce-api/stores/models"
	storeProviders "github.com/dangduoc08/ecommerce-api/stores/providers"
	userModels "github.com/dangduoc08/ecommerce-api/users/models"
	userProviders "github.com/dangduoc08/ecommerce-api/users/providers"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
)

type Seed struct {
	dbProviders.DB
	config.ConfigService
	common.Logger
	JWTAccessAPIKey    string
	JWTRefreshTokenKey string
	UserDBHandler      userProviders.DBHandler
	LocationDBHandler  locationProviders.DBHandler
	StoreDBHandler     storeProviders.DBHandler
	AuthCipher         authProviders.Cipher
}

func (instance Seed) NewProvider() core.Provider {

	// Create user_status enum
	instance.CreateEnum(constants.USER_STATUS_FIELD_NAME, constants.USER_STATUSES)

	// Create category_status enum
	instance.CreateEnum(constants.CATEGORY_STATUS_FIELD_NAME, constants.CATEGORY_STATUSES)

	// Create category_status enum
	instance.CreateEnum(constants.PRODUCT_STATUS_FIELD_NAME, constants.PRODUCT_STATUSES)

	// Create tables
	err := instance.AutoMigrate(&locationModels.Location{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&addressModels.Address{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&storeModels.Store{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&groupModels.Group{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&userModels.User{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&categoryModels.Category{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&manufacturerModels.Manufacturer{})
	if err != nil {
		panic(err)
	}

	err = instance.AutoMigrate(&productModels.Product{})
	if err != nil {
		panic(err)
	}

	// Seed locations
	totalLocations := instance.Count(constants.TABLE_LOCATION)
	if totalLocations == 0 {
		instance.LocationDBHandler.Seed(func(locationRec locationModels.Location) {
			resp := instance.Create(&locationRec)
			if resp.Error != nil {
				instance.Logger.Debug("SeedLocations", "error", resp.Error)
			}
		})
	}

	// Seed users
	hash, err := instance.AuthCipher.HashPassword(instance.ConfigService.Get("PASSWORD").(string))
	if err != nil {
		panic(err)
	}
	user := userModels.User{
		Username:  instance.ConfigService.Get("USERNAME").(string),
		Email:     instance.ConfigService.Get("EMAIL").(string),
		FirstName: instance.ConfigService.Get("FIRST_NAME").(string),
		LastName:  instance.ConfigService.Get("LAST_NAME").(string),
		Hash:      hash,
		Status:    userModels.UserStatus(constants.USER_ACTIVE),
		Groups: []*groupModels.Group{
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
		resp := instance.Create(&storeModels.Store{
			Name:      "Demo",
			Addresses: []addressModels.Address{},
			Users:     []userModels.User{user},
		})
		if resp.Error != nil {
			instance.Logger.Debug("SeedStores", "error", resp.Error)
		}
	}

	return instance
}
