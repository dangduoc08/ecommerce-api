package dtos

import (
	"strings"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type CREATE_Body_Data_DTO struct {
	StreetName string `bind:"street_name"`
	LocationID uint   `bind:"location_id"`
}

type CREATE_Body_DTO struct {
	Data CREATE_Body_Data_DTO `bind:"data"`
}

func (instance CREATE_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	bindedBody, _ := body.Bind(instance)
	bodyDTO := bindedBody.(CREATE_Body_DTO)

	bodyDTO.Data.StreetName = strings.TrimSpace(bodyDTO.Data.StreetName)

	return bodyDTO
}
