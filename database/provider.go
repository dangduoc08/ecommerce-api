package database

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/utils"
	"gorm.io/gorm"
)

type Provider struct {
	DB *gorm.DB
}

func (provider Provider) NewProvider() core.Provider {
	return provider
}

func (provider Provider) CreateType(typeName string, values []string) {
	formatValues := utils.ArrMap[string](values, func(el string, i int) string {
		return fmt.Sprintf("'%s'", el)
	})
	sql := fmt.Sprintf(`
		DO $$ BEGIN
			IF
				NOT EXISTS (SELECT oid FROM pg_type WHERE typname = '%v')
			THEN
				CREATE TYPE %v AS ENUM (%v);
			END IF;
		END $$;`, typeName, typeName, strings.Join(formatValues, ", "))

	resp := provider.DB.Exec(sql)
	if resp.Error != nil {
		panic(resp.Error)
	}
}
