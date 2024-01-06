package interceptors

import (
	"encoding/json"
	"reflect"

	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/exception"
)

type ChildCategoryTransformation struct{}

func (self ChildCategoryTransformation) Intercept(c gooh.Context, aggregation gooh.Aggregation) any {
	return aggregation.Pipe(
		aggregation.Consume(func(c gooh.Context, data any) any {
			if !reflect.ValueOf(data).IsNil() {
				if menu, ok := data.(*[]map[string]any); ok {
					for _, item := range *menu {
						if item["child_categories"] != "" {
							childCategories := []map[string]any{}

							err := json.Unmarshal([]byte(item["child_categories"].(string)), &childCategories)
							if err != nil {
								panic(exception.InternalServerErrorException(err.Error()))
							}

							childCategories = utils.ArrFilter(childCategories, func(el map[string]any, i int) bool {
								return el != nil
							})

							item["child_categories"] = childCategories
						}
					}
				}
			}

			return data
		}),
	)
}
