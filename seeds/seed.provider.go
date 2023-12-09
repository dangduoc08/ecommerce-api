package seeds

import (
	"github.com/dangduoc08/ecommerce-api/addresses"
	"github.com/dangduoc08/ecommerce-api/db"
	"github.com/dangduoc08/ecommerce-api/locations"
	"github.com/dangduoc08/ecommerce-api/stores"
	"github.com/dangduoc08/ecommerce-api/users"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
)

type SeedProvider struct {
	JWTAccessAPIKey    string
	JWTRefreshTokenKey string

	DBProvider       db.DBProvider
	ConfigService    config.ConfigService
	Logger           common.Logger
	UserProvider     users.UserProvider
	LocationProvider locations.LocationProvider
	StoreProvider    stores.StoreProvider
}

func (self SeedProvider) NewProvider() core.Provider {

	// Create user_status enum
	self.DBProvider.CreateEnum(users.STATUS, []string{
		users.ACTIVE,
		users.INACTIVE,
		users.SUSPENDED,
	})

	// // Create tables
	err := self.DBProvider.DB.AutoMigrate(&locations.Location{})
	if err != nil {
		panic(err)
	}

	err = self.DBProvider.DB.AutoMigrate(&addresses.Address{})
	if err != nil {
		panic(err)
	}

	err = self.DBProvider.DB.AutoMigrate(&stores.Store{})
	if err != nil {
		panic(err)
	}

	err = self.DBProvider.DB.AutoMigrate(&users.User{})
	if err != nil {
		panic(err)
	}

	// Seed locations
	totalLocations := self.DBProvider.Count(self.LocationProvider.GetModelName())
	if totalLocations == 0 {
		self.LocationProvider.Seed(func(locationRec locations.Location) {
			resp := self.DBProvider.DB.Create(&locationRec)
			if resp.Error != nil {
				self.Logger.Debug("SeedLocations", "error", resp.Error)
			}
		})
	}

	// Seed stores
	totalStores := self.DBProvider.Count(self.StoreProvider.GetModelName())
	if totalStores == 0 {
		var locationID1 uint = 1
		var locationID2 uint = 2
		resp := self.DBProvider.DB.Create(&stores.Store{
			Name:        "Demo",
			Description: "Demo shop",
			Email:       self.ConfigService.Get("EMAIL").(string),
			Addresses: []addresses.Address{
				{
					LocationID: &locationID1,
				},
				{
					LocationID: &locationID2,
				},
			},
		})
		if resp.Error != nil {
			self.Logger.Debug("SeedStores", "error", resp.Error)
		}
	}

	// Seed users
	totalUsers := self.DBProvider.Count(self.UserProvider.GetModelName())
	if totalUsers == 0 {
		hash, err := self.UserProvider.HashPassword(self.ConfigService.Get("PASSWORD").(string))
		if err != nil {
			panic(err)
		}
		resp := self.DBProvider.DB.Create(&users.User{
			Username:  self.ConfigService.Get("USERNAME").(string),
			Hash:      hash,
			Email:     self.ConfigService.Get("EMAIL").(string),
			FirstName: self.ConfigService.Get("FIRST_NAME").(string),
			LastName:  self.ConfigService.Get("LAST_NAME").(string),
			Status:    users.UserStatus(users.ACTIVE),
		})
		if resp.Error != nil {
			self.Logger.Debug("SeedUsers", "error", resp.Error)
		}
	}

	return self
}
