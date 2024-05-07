package seeds

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dangduoc08/ecommerce-api/admins/addresses"
	"github.com/dangduoc08/ecommerce-api/admins/auths"
	"github.com/dangduoc08/ecommerce-api/admins/categories"
	"github.com/dangduoc08/ecommerce-api/admins/groups"
	"github.com/dangduoc08/ecommerce-api/admins/locations"
	mailConfigurations "github.com/dangduoc08/ecommerce-api/admins/mail_configurations"
	"github.com/dangduoc08/ecommerce-api/admins/manufacturers"
	"github.com/dangduoc08/ecommerce-api/admins/products"
	"github.com/dangduoc08/ecommerce-api/admins/stores"
	"github.com/dangduoc08/ecommerce-api/admins/users"
	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/ecommerce-api/utils"
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
	if err := instance.AutoMigrate(&locations.LocationModel{}); err != nil {
		panic(err)
	}

	if err := instance.AutoMigrate(&addresses.AddressModel{}); err != nil {
		panic(err)
	}

	if err := instance.AutoMigrate(&stores.StoreModel{}); err != nil {
		panic(err)
	}

	if err := instance.AutoMigrate(&groups.GroupModel{}); err != nil {
		panic(err)
	}

	if err := instance.AutoMigrate(&users.UserModel{}); err != nil {
		panic(err)
	}

	if err := instance.AutoMigrate(&categories.CategoryModel{}); err != nil {
		panic(err)
	}

	if err := instance.AutoMigrate(&manufacturers.ManufacturerModel{}); err != nil {
		panic(err)
	}

	if err := instance.AutoMigrate(&products.ProductModel{}); err != nil {
		panic(err)
	}

	if err := instance.AutoMigrate(&mailConfigurations.MailConfigurationModel{}); err != nil {
		panic(err)
	}

	// Seed stores
	var totalStores int64
	instance.Model(&stores.StoreModel{}).Count(&totalStores)
	if totalStores == 0 {
		if resp := instance.Create(&stores.StoreModel{
			Name: "Demo",
		}); resp.Error != nil {
			instance.Logger.Debug("SeedStores", "error", resp.Error)
		}
	}
	storeRec := &stores.StoreModel{}
	instance.First(storeRec)

	// Seed locations
	var totalLocations int64
	instance.Model(&locations.LocationModel{}).Count(&totalLocations)
	if totalLocations == 0 {
		instance.SeedLocations(func(locationRec locations.LocationModel) {
			if resp := instance.Create(&locationRec); resp.Error != nil {
				instance.Logger.Debug("SeedLocations", "error", resp.Error)
			}
		})
	}

	// Seed users
	var totalUsers int64
	instance.Model(&users.UserModel{}).Count(&totalUsers)
	if totalUsers == 0 {
		instance.SeedUsers(storeRec.ID, func(userRec users.UserModel) {
			if resp := instance.Create(&userRec); resp.Error != nil {
				instance.Logger.Debug("SeedUsers", "error", resp.Error)
			}
		})
	}

	// Seed mail configuration
	var totalMailConfigurations int64
	instance.Model(&mailConfigurations.MailConfigurationModel{}).Count(&totalMailConfigurations)
	if totalMailConfigurations == 0 {
		if resp := instance.Create(&mailConfigurations.MailConfigurationModel{StoreID: storeRec.ID}); resp.Error != nil {
			instance.Logger.Debug("SeedMailConfiguration", "error", resp.Error)
		}
	}

	return instance
}

func (instance SeedProvider) SeedUsers(storeID uint, cb func(users.UserModel)) {
	dir, _ := os.Getwd()
	data, _ := os.ReadFile(fmt.Sprintf("%v/seeds/datas/%v", dir, "users.json"))

	var userList []map[string]any
	_ = json.Unmarshal(data, &userList)

	if len(userList) > 0 {
		for _, user := range userList {
			hash, _ := instance.AuthProvider.HashPassword(user["password"].(string))

			userRec := users.UserModel{
				Username:  user["user_name"].(string),
				Hash:      hash,
				Email:     user["email"].(string),
				FirstName: user["first_name"].(string),
				LastName:  user["last_name"].(string),
				Status:    users.UserStatus(constants.USER_ACTIVE),
				StoreID:   storeID,
				Groups: []*groups.GroupModel{
					{
						Name: user["groups"].(map[string]any)["name"].(string),
						Permissions: utils.ArrMap(
							user["groups"].(map[string]any)["permissions"].([]any),
							func(el any, idx int) string {
								return el.(string)
							}),
					},
				},
			}

			cb(userRec)
		}
	}
}

func (instance SeedProvider) SeedLocations(cb func(locations.LocationModel)) {
	dir, _ := os.Getwd()
	data, err := os.ReadFile(fmt.Sprintf("%v/seeds/datas/%v", dir, "locations.json"))
	if err != nil {
		panic(err)
	}

	var locationList []map[string]any
	err = json.Unmarshal(data, &locationList)
	if err != nil {
		panic(err)
	}

	if len(locationList) > 0 {
		for _, location := range locationList {
			locationRec := locations.LocationModel{
				LocationID: nil,
			}

			if f64ID, ok := location["id"].(float64); ok {
				locationRec.ID = uint(f64ID)
			}

			if f64LocationID, ok := location["location_id"].(float64); ok {
				unitLocationID := uint(f64LocationID)
				locationRec.LocationID = &unitLocationID
			}

			if strName, ok := location["name"].(string); ok {
				locationRec.Name = strName
			}

			if strSlug, ok := location["slug"].(string); ok {
				locationRec.Slug = strSlug
			}

			cb(locationRec)
		}
	}
}
