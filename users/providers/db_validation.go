package providers

import (
	"fmt"

	dbProviders "github.com/dangduoc08/ecommerce-api/db/providers"
	"github.com/dangduoc08/ecommerce-api/users/models"
	"github.com/dangduoc08/gooh/core"
)

type DBValidation struct {
	dbProviders.DB
}

func (self DBValidation) NewProvider() core.Provider {
	return self
}

func (self DBValidation) CheckDuplicated(data []map[string]string) bool {
	for _, kv := range data {
		for k, v := range kv {
			var userRec models.User
			self.Where(fmt.Sprintf("%v = ?", k), fmt.Sprintf("%v", v)).First(&userRec)
			if userRec.ID != 0 {
				return true
			}
		}
	}

	return false
}
