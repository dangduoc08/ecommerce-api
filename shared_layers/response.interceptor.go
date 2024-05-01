package sharedLayers

import (
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/modules/config"
)

type ResponseInterceptor struct {
	ConfigService config.ConfigService
}

func (instance ResponseInterceptor) Intercept(c gogo.Context, aggregation gogo.Aggregation) any {
	return aggregation.Pipe(
		aggregation.Consume(func(c gogo.Context, data any) any {
			transformedData := gogo.Map{
				"data": data,
			}
			return transformedData
		}),
	)
}
