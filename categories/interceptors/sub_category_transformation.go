package interceptors

import (
	"encoding/json"
	"reflect"

	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/exception"
)

type SubCategoryTransformation struct{}

func (instance SubCategoryTransformation) Intercept(c gooh.Context, aggregation gooh.Aggregation) any {
	return aggregation.Pipe(
		aggregation.Consume(func(c gooh.Context, data any) any {
			if !reflect.ValueOf(data).IsNil() {
				if menu, ok := data.(*[]map[string]any); ok {
					for _, item := range *menu {
						if item["sub_categories"] != "" {
							subCategories := []map[string]any{}

							err := json.Unmarshal([]byte(item["sub_categories"].(string)), &subCategories)
							if err != nil {
								panic(exception.InternalServerErrorException(err.Error()))
							}

							subCategories = utils.ArrFilter(subCategories, func(el map[string]any, i int) bool {
								return el != nil
							})

							item["categories"] = subCategories
						}
					}
				}
			}

			return data
		}),
	)
}
