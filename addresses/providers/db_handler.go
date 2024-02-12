package providers

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/addresses/models"
	"github.com/dangduoc08/ecommerce-api/constants"
	dbProviders "github.com/dangduoc08/ecommerce-api/dbs/providers"
	"github.com/dangduoc08/gooh/core"
)

type DBHandler struct {
	DBProvider dbProviders.DB
}

func (instance DBHandler) NewProvider() core.Provider {
	return instance
}

func (instance DBHandler) FindByID(id uint) (*models.Address, error) {
	addressRec := &models.Address{
		ID: id,
	}

	if err := instance.DBProvider.DB.
		Joins("Location.Location.Location").
		First(addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (instance DBHandler) FindOneBy(query *Query) (*models.Address, error) {
	addressRec := &models.Address{}
	addressQueries := map[string]any{}

	if query.ID != 0 {
		addressQueries["addresses.id"] = query.ID
	}

	if query.StoreID != 0 {
		addressQueries["store_id"] = query.StoreID
	}

	if err := instance.DBProvider.DB.
		Where(addressQueries).
		Joins("Location.Location.Location").
		First(addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (instance DBHandler) FindManyBy(query *Query) ([]*models.Address, error) {
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

	if err := instance.DBProvider.DB.
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

func (instance DBHandler) CreateOne(data *Creation) (*models.Address, error) {
	addressRec := &models.Address{
		StoreID:    data.StoreID,
		StreetName: data.StreetName,
		LocationID: data.LocationID,
	}

	if err := instance.DBProvider.DB.Create(&addressRec).Error; err != nil {
		return nil, err
	}

	if err := instance.DBProvider.DB.
		Joins("Location.Location.Location").
		First(&addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (instance DBHandler) UpdateByID(id uint, data *Update) (*models.Address, error) {
	addressRec := &models.Address{
		ID:         id,
		StoreID:    data.StoreID,
		StreetName: data.StreetName,
		LocationID: data.LocationID,
	}

	if err := instance.DBProvider.DB.Updates(&addressRec).Error; err != nil {
		return nil, err
	}

	if err := instance.DBProvider.DB.
		Joins("Location.Location.Location").
		First(&addressRec).
		Error; err != nil {
		return nil, err
	}

	return addressRec, nil
}

func (instance DBHandler) DeleteByID(id uint) error {
	return instance.DBProvider.DB.Delete(&models.Address{}, id).Error
}
