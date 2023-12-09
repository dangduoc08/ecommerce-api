package addresses

import (
	"path"
	"reflect"

	"github.com/dangduoc08/ecommerce-api/db"
	"github.com/dangduoc08/ecommerce-api/locations"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
)

type AddressProvider struct {
	DBProvider       db.DBProvider
	ConfigService    config.ConfigService
	LocationProvider locations.LocationProvider
}

func (self AddressProvider) NewProvider() core.Provider {
	return self
}

func (self AddressProvider) GetModelName() string {
	return path.Base(reflect.TypeOf(self).PkgPath())
}

func (self AddressProvider) FindManyBy(query *AddressQuery) ([]Address, error) {
	addressRecs := []Address{}

	addressQueries := map[string]any{}

	if query.ID != 0 {
		addressQueries["addresses.id"] = query.ID
	}

	if query.StoreID != 0 {
		addressQueries["store_id"] = query.StoreID
	}

	resp := self.DBProvider.DB.
		Joins("Location.Location.Location").
		Where(addressQueries).
		Find(&addressRecs)

	if resp.Error != nil {
		return []Address{}, resp.Error
	}

	return addressRecs, nil
}

func (self AddressProvider) FindByID(id uint) (*Address, error) {
	addressRec := &Address{
		ID: id,
	}

	resp := self.DBProvider.DB.First(addressRec)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return addressRec, nil
}

func (self AddressProvider) CreateOne(data *AddressCreation) (*Address, error) {
	addressRec := &Address{
		StoreID:    data.StoreID,
		StreetName: data.StreetName,
		LocationID: data.LocationID,
	}

	resp := self.DBProvider.DB.Create(&addressRec)

	if resp.Error != nil {
		return nil, resp.Error
	}

	self.DBProvider.DB.Joins("Location.Location.Location").First(&addressRec)

	return addressRec, resp.Error
}

func (self AddressProvider) UpdateByID(id uint, data *AddressUpdate) (*Address, error) {
	addressRec := &Address{
		ID:         id,
		StoreID:    data.StoreID,
		StreetName: data.StreetName,
		LocationID: data.LocationID,
	}

	resp := self.DBProvider.DB.Save(&addressRec)

	if resp.Error != nil {
		return nil, resp.Error
	}

	self.DBProvider.DB.Joins("Location.Location.Location").First(&addressRec)

	return addressRec, resp.Error
}
