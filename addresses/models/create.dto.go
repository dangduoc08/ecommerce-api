package models

import (
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
)

type CREATE_Body_Data struct {
	StreetName string `bind:"street_name"`
	LocationID uint   `bind:"location_id"`
}

type CREATE_Body struct {
	Data CREATE_Body_Data `bind:"data"`
}

func (self CREATE_Body) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	bodyDTO, _ := body.Bind(self)
	return bodyDTO
}
