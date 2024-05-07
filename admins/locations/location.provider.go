package locations

import (
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/modules/config"
)

type LocationProvider struct {
	dbs.DBProvider
	config.ConfigService
}

func (instance LocationProvider) NewProvider() core.Provider {
	return instance
}

func (instance LocationProvider) FindByID(id uint) (*LocationModel, error) {
	locationRec := &LocationModel{
		ID: id,
	}

	resp := instance.First(locationRec)

	return locationRec, resp.Error
}

func (instance LocationProvider) FindManyBy(query *Query) ([]*LocationModel, error) {
	var locations []*LocationModel

	resp := instance.
		Select("id", "location_id", "name", "slug").
		Where("location_id", query.LocationID).
		Find(&locations)

	return locations, resp.Error
}
