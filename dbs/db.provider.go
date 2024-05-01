package dbs

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/utils"
	"gorm.io/gorm"
)

type DBProvider struct {
	*gorm.DB
}

func (instance DBProvider) NewProvider() core.Provider {

	return instance
}

func (instance DBProvider) CreateEnum(typeName string, values []string) {
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

	resp := instance.DB.Exec(sql)
	if resp.Error != nil {
		panic(resp.Error)
	}
}

func (instance DBProvider) Count(tableName string) int {
	var count int
	if err := instance.DB.Raw(fmt.Sprintf("SELECT count(*) FROM %v", tableName)).Scan(&count).Error; err != nil {
		panic(err)
	}

	return count
}
