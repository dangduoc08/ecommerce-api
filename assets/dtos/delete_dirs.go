package dtos

import (
	"github.com/dangduoc08/ecommerce-api/assets/providers"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
)

type DELETE_Query struct {
	providers.HandleAsset
	Dirs []string `bind:"dirs"`
}

func (instance DELETE_Query) Transform(query gogo.Query, medata common.ArgumentMetadata) any {
	bindedStruct, _ := query.Bind(instance)

	queryDTO := bindedStruct.(DELETE_Query)

	for i, dir := range queryDTO.Dirs {
		queryDTO.Dirs[i] = instance.CleanDir(dir)
	}
	queryDTO.Dirs = utils.ArrToUnique(queryDTO.Dirs)

	return queryDTO
}
