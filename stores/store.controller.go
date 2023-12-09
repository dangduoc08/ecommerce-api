package stores

import (
	"strconv"

	"github.com/dangduoc08/ecommerce-api/globals"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/exception"
	"github.com/dangduoc08/gooh/modules/config"
)

type StoreController struct {
	common.Rest
	common.Guard
	Logger        common.Logger
	StoreProvider StoreProvider
	ConfigService config.ConfigService
}

func (self StoreController) NewController() core.Controller {
	self.
		Prefix("v1").
		Prefix("stores")

	self.
		BindGuard(globals.AccessAPIGuard{})

	return self
}

func (self StoreController) UPDATE_BY_id() {

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

	store, err := self.StoreProvider.GetOneByID(uint(id))
	if err != nil {
		panic(exception.NotFoundException(err.Error()))
	}

	return store
}
