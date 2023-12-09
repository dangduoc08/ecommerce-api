package db

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/conf"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBModule = func(
	logger common.Logger,
) *core.Module {

	sslmode := "disable"
	if conf.Value.PostgresSSL {
		sslmode = "require"
	}

	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=Asia/Shanghai",
		conf.Value.PostgresHost,
		conf.Value.PostgresUser,
		conf.Value.PostgresPassword,
		conf.Value.PostgresDB,
		conf.Value.PostgresPort,
		sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})

	if err != nil {
		logger.Fatal("PostgresSQL", "error", err.Error(), "connected", false)
	} else {
		logger.Info("PostgresSQL", "connected", true)
	}

	provider := DBProvider{
		DB: db,
	}

	module := core.ModuleBuilder().
		Exports(provider).
		Build()

	module.IsGlobal = true

	return module
}
