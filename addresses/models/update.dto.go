package models

import (
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
)

type UPDATE_BY_id_Param struct {
	ID uint `bind:"id"`
}

func (self UPDATE_BY_id_Param) Transform(param gooh.Param, medata common.ArgumentMetadata) any {
	paramDTO, _ := param.Bind(self)
	return paramDTO
}

type UPDATE_BY_id_Body_Data struct {
	StreetName string `bind:"street_name"`
	LocationID uint   `bind:"location_id"`
}

type UPDATE_BY_id_Body struct {
	Data UPDATE_BY_id_Body_Data `bind:"data"`
}

func (self UPDATE_BY_id_Body) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	bodyDTO, _ := body.Bind(self)
	return bodyDTO
}
