package providers

import (
	addressModels "github.com/dangduoc08/ecommerce-api/addresses/models"
	authProviders "github.com/dangduoc08/ecommerce-api/auths/providers"
	categoryModels "github.com/dangduoc08/ecommerce-api/categories/models"
	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	groupModels "github.com/dangduoc08/ecommerce-api/groups/models"
	locationModels "github.com/dangduoc08/ecommerce-api/locations/models"
	locationProviders "github.com/dangduoc08/ecommerce-api/locations/providers"
	manufacturerModels "github.com/dangduoc08/ecommerce-api/manufacturers/models"
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

func (self Seed) NewProvider() core.Provider {

	// Create user_status enum
	self.CreateEnum(constants.USER_STATUS_FIELD_NAME, constants.USER_STATUSES)

	// Create category_status enum
	self.CreateEnum(constants.CATEGORY_STATUS_FIELD_NAME, constants.CATEGORY_STATUSES)

	// Create tables
	err := self.AutoMigrate(&locationModels.Location{})
	if err != nil {
		panic(err)
	}

	err = self.AutoMigrate(&addressModels.Address{})
	if err != nil {
		panic(err)
	}

	err = self.AutoMigrate(&storeModels.Store{})
	if err != nil {
		panic(err)
	}

	err = self.AutoMigrate(&groupModels.Group{})
	if err != nil {
		panic(err)
	}

	err = self.AutoMigrate(&userModels.User{})
	if err != nil {
		panic(err)
	}

	err = self.AutoMigrate(&categoryModels.Category{})
	if err != nil {
		panic(err)
	}

	err = self.AutoMigrate(&manufacturerModels.Manufacturer{})
	if err != nil {
		panic(err)
	}

	// Seed locations
	totalLocations := self.Count(constants.TABLE_LOCATION)
	if totalLocations == 0 {
		self.LocationDBHandler.Seed(func(locationRec locationModels.Location) {
			resp := self.Create(&locationRec)
			if resp.Error != nil {
				self.Logger.Debug("SeedLocations", "error", resp.Error)
			}
		})
	}

	// Seed users
	hash, err := self.AuthCipher.HashPassword(self.ConfigService.Get("PASSWORD").(string))
	if err != nil {
		panic(err)
	}
	user := userModels.User{
		Username:  self.ConfigService.Get("USERNAME").(string),
		Email:     self.ConfigService.Get("EMAIL").(string),
		FirstName: self.ConfigService.Get("FIRST_NAME").(string),
		LastName:  self.ConfigService.Get("LAST_NAME").(string),
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
	totalStores := self.Count(constants.TABLE_STORE)
	if totalStores == 0 {
		resp := self.Create(&storeModels.Store{
			Name:      "Demo",
			Addresses: []addressModels.Address{},
			Users:     []userModels.User{user},
		})
		if resp.Error != nil {
			self.Logger.Debug("SeedStores", "error", resp.Error)
		}
	}

	return self
}
