package addresses

import (
	"github.com/dangduoc08/ecommerce-api/locations"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/context"
	"github.com/dangduoc08/gooh/utils"
)

type AddressInterceptor struct{}

func (self AddressInterceptor) transformLocationObjToArr(location *locations.Location, locationArr []*locations.Location) []*locations.Location {

	locationArr = append(
		[]*locations.Location{
			{
				ID:   location.ID,
				Name: location.Name,
				Slug: location.Slug,
			},
		},
		locationArr...,
	)

	if location.Location != nil {
		return self.transformLocationObjToArr(location.Location, locationArr)
	}

	return locationArr
}

func (self AddressInterceptor) toResponseObj(address *Address) gooh.Map {
	respData := gooh.Map{
		"id":          address.ID,
		"street_name": address.StreetName,
	}

	if address.Location != nil {
		locationArr := self.transformLocationObjToArr(address.Location, []*locations.Location{})
		respData["locations"] = locationArr
	} else {
		respData["locations"] = []any{}
	}

	return respData
}

func (self AddressInterceptor) Intercept(c gooh.Context, aggregation gooh.Aggregation) any {
	return aggregation.Pipe(
		aggregation.Consume(func(ctx gooh.Context, data any) any {
			switch address := data.(type) {

			case *Address:
				return self.toResponseObj(address)

			case []Address:
				return utils.ArrMap[Address, gooh.Map](address, func(address Address, i int) context.Map {
					return self.toResponseObj(&address)
				})
			}

			return data
		}),
	)
}
