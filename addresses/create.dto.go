package addresses

import (
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
)

type CREATE_Body_Data_DTO struct {
	StreetName string `bind:"street_name"`
	LocationID uint   `bind:"location_id"`
}

type CREATE_Body_DTO struct {
	Data CREATE_Body_Data_DTO `bind:"data"`
}

func (self CREATE_Body_DTO) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	return body.Bind(self)
}
