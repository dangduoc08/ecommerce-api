package dbs

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	conf "github.com/dangduoc08/ecommerce-api/confs"
	"github.com/dangduoc08/ecommerce-api/dbs/providers"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	common.Logger
}

func (instance *GormLogger) LogMode(logLevel gormLogger.LogLevel) gormLogger.Interface {
	return instance
}

func (instance *GormLogger) Info(c context.Context, msg string, data ...any) {
	instance.Logger.Info(msg, data)
}

func (instance *GormLogger) Warn(c context.Context, msg string, data ...any) {
	instance.Logger.Warn(msg, data)
}

func (instance *GormLogger) Error(c context.Context, msg string, data ...any) {
	instance.Logger.Error(msg, data)
}

func (instance *GormLogger) Trace(c context.Context, begin time.Time, cb func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := cb()
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")
	sql = strings.TrimSpace(sql)

	if err != nil {
		instance.Logger.Error(err.Error(), "sql", sql, "rowsAffected", rowsAffected)
	} else {
		instance.Logger.Debug("GORM", "sql", sql, "rowsAffected", rowsAffected)
	}
}

var Module = func(
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
		Logger: &GormLogger{logger},
	})

	if err != nil {
		logger.Fatal("PostgresSQL", "error", err.Error(), "connected", false)
	} else {
		logger.Info("PostgresSQL", "connected", true)
	}

	provider := providers.DB{
		DB: db,
	}

	module := core.ModuleBuilder().
		Providers(provider).
		Build()

	module.IsGlobal = true

	return module
}
