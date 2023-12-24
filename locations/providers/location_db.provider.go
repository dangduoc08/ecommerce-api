package providers

import (
	"encoding/json"
	"fmt"
	"os"

	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	"github.com/dangduoc08/ecommerce-api/locations/models"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
)

type LocationDB struct {
	DBProvider    dbProviders.DB
	ConfigService config.ConfigService
}

func (self LocationDB) NewProvider() core.Provider {
	return self
}

func (self LocationDB) Seed(cb func(models.Location)) {
	dir, _ := os.Getwd()
	data, err := os.ReadFile(fmt.Sprintf("%v/seeds/data/%v", dir, "locations.json"))
	var locationList []map[string]any
	err = json.Unmarshal(data, &locationList)
	if err != nil {
		panic(err)
	}
	if len(locationList) > 0 {
		for _, location := range locationList {
			locationRec := models.Location{
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

func (self LocationDB) FindByID(id uint) (*models.Location, error) {
	locationRec := &models.Location{
		ID: id,
	}

	resp := self.DBProvider.DB.First(locationRec)

	return locationRec, resp.Error
}

func (self LocationDB) FindManyBy(query *LocationQuery) ([]*models.Location, error) {
	var locations []*models.Location

	resp := self.DBProvider.DB.
		Select("id", "location_id", "name", "slug").
		Where("location_id", query.LocationID).
		Find(&locations)

	return locations, resp.Error
}
