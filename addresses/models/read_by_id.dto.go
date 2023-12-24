package models

import (
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
)

type READ_BY_id_Param struct {
	ID uint `bind:"id"`
}

func (self READ_BY_id_Param) Transform(param gooh.Param, medata common.ArgumentMetadata) any {
	paramDTO, _ := param.Bind(self)
	return paramDTO
}
