package addresses

import (
	"reflect"

	"github.com/dangduoc08/ecommerce-api/admins/locations"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/ctx"
	"github.com/dangduoc08/gogo/utils"
)

type AddressInterceptor struct{}

func (instance AddressInterceptor) transformLocationObjToArr(
	location *locations.LocationModel,
	locationArr []*locations.LocationModel,
) []*locations.LocationModel {
	locationArr = append(
		[]*locations.LocationModel{
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

func (instance AddressInterceptor) toResponseObj(address *AddressModel) gogo.Map {
	respData := gogo.Map{
		"id":          address.ID,
		"street_name": address.StreetName,
		"created_at":  address.CreatedAt,
		"updated_at":  address.UpdatedAt,
	}

	if address.Location != nil {
		locationArr := instance.transformLocationObjToArr(address.Location, []*locations.LocationModel{})
		respData["locations"] = locationArr
	} else {
		respData["locations"] = []any{}
	}

	return respData
}

func (instance AddressInterceptor) Intercept(c gogo.Context, aggregation gogo.Aggregation) any {
	return aggregation.Pipe(
		aggregation.Consume(func(c gogo.Context, data any) any {
			if !reflect.ValueOf(data).IsNil() {
				switch address := data.(type) {

				case *AddressModel:
					return instance.toResponseObj(address)

				case []*AddressModel:
					return utils.ArrMap[*AddressModel, gogo.Map](address, func(address *AddressModel, i int) ctx.Map {
						return instance.toResponseObj(address)
					})
				}
			}

			return data
		}),
	)
}
