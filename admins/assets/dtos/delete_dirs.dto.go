package dtos

import (
	"github.com/dangduoc08/ecommerce-api/admins/assets/commons"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type DELETE_Query_DTO struct {
	CommonProvider commons.CommonProvider
	Dirs           []string `bind:"dirs"`
}

func (instance DELETE_Query_DTO) Transform(query gogo.Query, medata common.ArgumentMetadata) any {
	bindedStruct, _ := query.Bind(instance)

	queryDTO := bindedStruct.(DELETE_Query_DTO)

	for i, dir := range queryDTO.Dirs {
		queryDTO.Dirs[i] = instance.CommonProvider.CleanDir(dir)
	}
	queryDTO.Dirs = utils.ArrToUnique(queryDTO.Dirs)

	return queryDTO
}
