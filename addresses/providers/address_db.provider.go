package providers

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/addresses/models"
	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	"github.com/dangduoc08/gooh/core"
)

type AddressDB struct {
	DBProvider dbProviders.DB
}

func (self AddressDB) NewProvider() core.Provider {
	return self
}

func (self AddressDB) FindByID(id uint) (*models.Address, error) {
	addressRec := &models.Address{
		ID: id,
	}

	if err := self.DBProvider.DB.
		Joins("Location.Location.Location").
		First(addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (self AddressDB) FindOneBy(query *AddressQuery) (*models.Address, error) {
	addressRec := &models.Address{}
	addressQueries := map[string]any{}

	if query.ID != 0 {
		addressQueries["addresses.id"] = query.ID
	}

	if query.StoreID != 0 {
		addressQueries["store_id"] = query.StoreID
	}

	if err := self.DBProvider.DB.
		Where(addressQueries).
		Joins("Location.Location.Location").
		First(addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (self AddressDB) FindManyBy(query *AddressQuery) ([]*models.Address, error) {
	addressRecs := []*models.Address{}
	addressQueries := map[string]any{}

	if query.ID != 0 {
		addressQueries["addresses.id"] = query.ID
	}

	if query.StoreID != 0 {
		addressQueries["store_id"] = query.StoreID
	}

	if query.Order == "" {
		query.Order = constants.ASC
	}

	sort := ""
	if query.Sort != "" {
		sort = fmt.Sprintf("%v %v", query.Sort, query.Order)
	}

	if err := self.DBProvider.DB.
		Order(sort).
		Limit(query.Limit).
		Offset(query.Offset).
		Joins("Location.Location.Location").
		Where(addressQueries).
		Find(&addressRecs).
		Error; err != nil {
		return []*models.Address{}, err
	}

	return addressRecs, nil
}

func (self AddressDB) CreateOne(data *AddressCreation) (*models.Address, error) {
	addressRec := &models.Address{
		StoreID:    data.StoreID,
		StreetName: data.StreetName,
		LocationID: data.LocationID,
	}

	if err := self.DBProvider.DB.Create(&addressRec).Error; err != nil {
		return nil, err
	}

	if err := self.DBProvider.DB.
		Joins("Location.Location.Location").
		First(&addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (self AddressDB) UpdateByID(id uint, data *AddressUpdate) (*models.Address, error) {
	addressRec := &models.Address{
		ID:         id,
		StoreID:    data.StoreID,
		StreetName: data.StreetName,
		LocationID: data.LocationID,
	}

	if err := self.DBProvider.DB.Updates(&addressRec).Error; err != nil {
		return nil, err
	}

	if err := self.DBProvider.DB.
		Joins("Location.Location.Location").
		First(&addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (self AddressDB) DeleteByID(id uint) error {
	return self.DBProvider.DB.Delete(&models.Address{}, id).Error
}
