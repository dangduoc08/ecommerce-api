package dtos

import (
	"strings"

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
	bindedBody, _ := body.Bind(self)
	bodyDTO := bindedBody.(CREATE_Body)

	bodyDTO.Data.StreetName = strings.TrimSpace(bodyDTO.Data.StreetName)

	return bodyDTO
}
