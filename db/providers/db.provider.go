package providers

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/utils"
	"gorm.io/gorm"
)

type DB struct {
	DB *gorm.DB
}

func (dbProvider DB) NewProvider() core.Provider {

	return dbProvider
}

func (dbProvider DB) CreateEnum(typeName string, values []string) {
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

	resp := dbProvider.DB.Exec(sql)
	if resp.Error != nil {
		panic(resp.Error)
	}
}

func (formatValues DB) Count(tableName string) int {
	var count int
	if err := formatValues.DB.Raw(fmt.Sprintf("SELECT count(*) FROM %v", tableName)).Scan(&count).Error; err != nil {
		panic(err)
	}

	return count
}
