package database

import (
	"fmt"

	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/modules/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Module = func(
	configService config.ConfigService,
	logger common.Logger,
) *core.Module {

	sslmode := "disable"
	if configService.Get("POSTGRES_SSL").(bool) {
		sslmode = "require"
	}

	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=Asia/Shanghai",
		configService.Get("POSTGRES_HOST"),
		configService.Get("POSTGRES_USER"),
		configService.Get("POSTGRES_PASSWORD"),
		configService.Get("POSTGRES_DB"),
		configService.Get("POSTGRES_PORT"),
		sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Fatal("PostgresSQL", "error", err.Error(), "connected", false)
	} else {
		logger.Info("PostgresSQL", "connected", true)
	}

	provider := Provider{
		DB: db,
	}

	module := core.ModuleBuilder().
		Exports(provider).
		Build()

	return module
}
