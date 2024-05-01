package addresses

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/gogo/core"
)

type AddressProvider struct {
	dbs.DBProvider
}

func (instance AddressProvider) NewProvider() core.Provider {
	return instance
}

func (instance AddressProvider) FindByID(id uint) (*AddressModel, error) {
	addressRec := &AddressModel{
		ID: id,
	}

	if err := instance.DB.
		Joins("Location.Location.Location").
		First(addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (instance AddressProvider) FindOneBy(query *Query) (*AddressModel, error) {
	addressRec := &AddressModel{}
	addressQueries := map[string]any{}

	if query.ID != 0 {
		addressQueries["addresses.id"] = query.ID
	}

	if query.StoreID != 0 {
		addressQueries["store_id"] = query.StoreID
	}

	if err := instance.DB.
		Where(addressQueries).
		Joins("Location.Location.Location").
		First(addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (instance AddressProvider) FindManyBy(query *Query) ([]*AddressModel, error) {
	addressRecs := []*AddressModel{}
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

	if err := instance.DB.
		Order(sort).
		Limit(query.Limit).
		Offset(query.Offset).
		Joins("Location.Location.Location").
		Where(addressQueries).
		Find(&addressRecs).
		Error; err != nil {
		return []*AddressModel{}, err
	}

	return addressRecs, nil
}

func (instance AddressProvider) CreateOne(data *Creation) (*AddressModel, error) {
	addressRec := &AddressModel{
		StoreID:    data.StoreID,
		StreetName: data.StreetName,
		LocationID: data.LocationID,
	}

	if err := instance.DB.Create(&addressRec).Error; err != nil {
		return nil, err
	}

	if err := instance.DB.
		Joins("Location.Location.Location").
		First(&addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (instance AddressProvider) UpdateByID(id uint, data *Update) (*AddressModel, error) {
	addressRec := &AddressModel{
		ID:         id,
		StoreID:    data.StoreID,
		StreetName: data.StreetName,
		LocationID: data.LocationID,
	}

	if err := instance.DB.Updates(&addressRec).Error; err != nil {
		return nil, err
	}

	if err := instance.DB.
		Joins("Location.Location.Location").
		First(&addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (instance AddressProvider) DeleteByID(id uint) error {
	return instance.DB.Delete(&AddressModel{}, id).Error
}
