package statics

import (
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
)

type StaticController struct {
	common.REST
	PublicPath string
}

func (instance StaticController) NewController() core.Controller {
	return instance
}

func (instance StaticController) SERVE_ANY_VERSION_NEUTRAL() string {
	return instance.PublicPath
}
