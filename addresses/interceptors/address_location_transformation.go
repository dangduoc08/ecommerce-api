package interceptors

import (
	"reflect"

	"github.com/dangduoc08/ecommerce-api/addresses/models"
	locationModels "github.com/dangduoc08/ecommerce-api/locations/models"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/ctx"
	"github.com/dangduoc08/gogo/utils"
)

type AddressLocationTransformation struct{}

func (instance AddressLocationTransformation) transformLocationObjToArr(
	location *locationModels.Location,
	locationArr []*locationModels.Location,
) []*locationModels.Location {
	locationArr = append(
		[]*locationModels.Location{
			{
				ID:   location.ID,
				Name: location.Name,
				Slug: location.Slug,
			},
		},
		locationArr...,
	)

	if location.Location != nil {
		return instance.transformLocationObjToArr(location.Location, locationArr)
	}

	return locationArr
}

func (instance AddressLocationTransformation) toResponseObj(address *models.Address) gogo.Map {
	respData := gogo.Map{
		"id":          address.ID,
		"street_name": address.StreetName,
		"created_at":  address.CreatedAt,
		"updated_at":  address.UpdatedAt,
	}

	if address.Location != nil {
		locationArr := instance.transformLocationObjToArr(address.Location, []*locationModels.Location{})
		respData["locations"] = locationArr
	} else {
		respData["locations"] = []any{}
	}

	return respData
}

func (instance AddressLocationTransformation) Intercept(c gogo.Context, aggregation gogo.Aggregation) any {
	return aggregation.Pipe(
		aggregation.Consume(func(c gogo.Context, data any) any {
			if !reflect.ValueOf(data).IsNil() {
				switch address := data.(type) {

				case *models.Address:
					return instance.toResponseObj(address)

				case []*models.Address:
					return utils.ArrMap[*models.Address, gogo.Map](address, func(address *models.Address, i int) ctx.Map {
						return instance.toResponseObj(address)
					})
				}
			}

			return data
		}),
	)
}
