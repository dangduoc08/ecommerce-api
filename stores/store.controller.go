package stores

import (
	"fmt"
	"strconv"

	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
)

type StoreController struct {
	common.REST
	common.Guard
	Logger        common.Logger
	StoreProvider StoreProvider
	ConfigService config.ConfigService
}

func (self StoreController) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("stores")

	return self
}

func (self StoreController) CREATE_addresses_OF_BY_id() {

}

func (self StoreController) UPDATE_BY_id(
	param gooh.Param,
	body UPDATE_BY_ID_Body_DTO,
) {
	var storeID uint

	if param.Get("id") != "" {
		u64StoreID, err := strconv.ParseUint(param.Get("id"), 10, 64)
		if err != nil {
			panic(exception.UnprocessableEntityException(err.Error()))
		}
		storeID = uint(u64StoreID)
	}

	fmt.Println(storeID, body)
}

func (self StoreController) READ_BY_id(
	param gooh.Param,
) *Store {
	storeID := param.Get("id")
	id, err := strconv.ParseUint(storeID, 10, 0)
	if err != nil {
		self.Logger.Error(
			"strconv.ParseUint",
			"param.Get(\"id\")", param.Get("id"),
			"errMsg", err,
		)
		panic(exception.UnprocessableEntityException(err.Error()))
	}

	store, err := self.StoreProvider.FindOneByID(uint(id))
	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	return store
}
