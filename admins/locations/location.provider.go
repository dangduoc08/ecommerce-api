package locations

import (
	"encoding/json"
	"fmt"
	"os"

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

func (instance LocationProvider) Seed(cb func(LocationModel)) {
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
			locationRec := LocationModel{
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
