package dtos

import (
	"github.com/dangduoc08/ecommerce-api/assets/providers"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
)

type DELETE_Query struct {
	providers.HandleAsset
	Dirs []string `bind:"dirs"`
}

func (self DELETE_Query) Transform(query gooh.Query, medata common.ArgumentMetadata) any {
	bindedStruct, _ := query.Bind(self)

	queryDTO := bindedStruct.(DELETE_Query)

	for i, dir := range queryDTO.Dirs {
		queryDTO.Dirs[i] = self.CleanDir(dir)
	}
	queryDTO.Dirs = utils.ArrToUnique(queryDTO.Dirs)

	return queryDTO
}
