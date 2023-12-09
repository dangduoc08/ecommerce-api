package locations

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/dangduoc08/ecommerce-api/db"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
)

type LocationProvider struct {
	DBProvider    db.DBProvider
	ConfigService config.ConfigService
}

func (self LocationProvider) NewProvider() core.Provider {
	return self
}

func (self LocationProvider) GetModelName() string {
	return path.Base(reflect.TypeOf(self).PkgPath())
}

func (self LocationProvider) Seed(cb func(Location)) {
	dir, _ := os.Getwd()
	data, err := os.ReadFile(fmt.Sprintf("%v/seeds/data/%v", dir, "locations.json"))
	var locationList []map[string]any
	err = json.Unmarshal(data, &locationList)
	if err != nil {
		panic(err)
	}
	if len(locationList) > 0 {
		for _, location := range locationList {
			locationRec := Location{
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

func (self LocationProvider) FindAllBy(locationID *uint) ([]Location, error) {
	var locations []Location

	resp := self.DBProvider.DB.
		Select("id", "location_id", "name", "slug").
		Where("location_id", locationID).
		Find(&locations)

	return locations, resp.Error
}
