package providers

import (
	addressModels "github.com/dangduoc08/ecommerce-api/addresses/models"
	authProviders "github.com/dangduoc08/ecommerce-api/auths/providers"
	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	groupModels "github.com/dangduoc08/ecommerce-api/groups/models"
	locationModels "github.com/dangduoc08/ecommerce-api/locations/models"
	locationProviders "github.com/dangduoc08/ecommerce-api/locations/providers"
	storeModels "github.com/dangduoc08/ecommerce-api/stores/models"
	storeProviders "github.com/dangduoc08/ecommerce-api/stores/providers"
	userModels "github.com/dangduoc08/ecommerce-api/users/models"
	userProviders "github.com/dangduoc08/ecommerce-api/users/providers"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
)

type Seed struct {
	JWTAccessAPIKey    string
	JWTRefreshTokenKey string

	DBProvider    dbProviders.DB
	ConfigService config.ConfigService
	Logger        common.Logger
	UserDB        userProviders.UserDB
	LocationDB    locationProviders.LocationDB
	StoreDB       storeProviders.StoreDB
	AuthHandler   authProviders.AuthHandler
}

func (self Seed) NewProvider() core.Provider {

	// Create user_status enum
	self.DBProvider.CreateEnum(constants.USER_STATUS_FIELD_NAME, constants.USER_STATUSES)

	// Create tables
	err := self.DBProvider.DB.AutoMigrate(&locationModels.Location{})
	if err != nil {
		panic(err)
	}

	err = self.DBProvider.DB.AutoMigrate(&addressModels.Address{})
	if err != nil {
		panic(err)
	}

	err = self.DBProvider.DB.AutoMigrate(&storeModels.Store{})
	if err != nil {
		panic(err)
	}

	err = self.DBProvider.DB.AutoMigrate(&groupModels.Group{})
	if err != nil {
		panic(err)
	}

	err = self.DBProvider.DB.AutoMigrate(&userModels.User{})
	if err != nil {
		panic(err)
	}

	// Seed locations
	totalLocations := self.DBProvider.Count(constants.TABLE_LOCATION)
	if totalLocations == 0 {
		self.LocationDB.Seed(func(locationRec locationModels.Location) {
			resp := self.DBProvider.DB.Create(&locationRec)
			if resp.Error != nil {
				self.Logger.Debug("SeedLocations", "error", resp.Error)
			}
		})
	}

	// Seed users
	hash, err := self.AuthHandler.HashPassword(self.ConfigService.Get("PASSWORD").(string))
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
	totalStores := self.DBProvider.Count(constants.TABLE_STORE)
	if totalStores == 0 {
		resp := self.DBProvider.DB.Create(&storeModels.Store{
			Name:        "Demo",
			Description: "Demo shop",
			Email:       self.ConfigService.Get("EMAIL").(string),
			Addresses:   []addressModels.Address{},
			Users:       []userModels.User{user},
		})
		if resp.Error != nil {
			self.Logger.Debug("SeedStores", "error", resp.Error)
		}
	}

	return self
}
