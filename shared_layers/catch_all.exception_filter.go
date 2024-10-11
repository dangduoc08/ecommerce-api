package sharedLayers

import (
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/exception"
)

type AllExceptionsFilter struct{}

func (allExceptionsFilter AllExceptionsFilter) Catch(c gogo.Context, ex gogo.Exception) {
	errResp := ex.GetResponse()
	locale := c.Header().Get("locale")
	if _, ok := utils.Translation[locale]; !ok || locale == "" {
		locale = "en_US"
	}

	switch e := errResp.(type) {
	case []string:
		errResp = utils.Reason(locale, e...)
	case string:
		errResp = utils.Reason(locale, e)
	case []map[string]any:
		for _, m := range e {
			for k, v := range m {
				if k == "reason" {
					switch v := v.(type) {
					case []string:
						m[k] = utils.Reason(locale, v...)
					case string:
						m[k] = utils.Reason(locale, v)
					}
				}
			}
		}
	}

	panic(exception.NewException(errResp, ex.GetCode(), ex))
}
