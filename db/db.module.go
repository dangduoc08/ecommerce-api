package db

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dangduoc08/ecommerce-api/conf"
	"github.com/dangduoc08/ecommerce-api/db/providers"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	common.Logger
}

func (self *GormLogger) LogMode(logLevel gormLogger.LogLevel) gormLogger.Interface {
	return self
}

func (self *GormLogger) Info(c context.Context, msg string, data ...any) {
	self.Logger.Info(msg, data)
}

func (self *GormLogger) Warn(c context.Context, msg string, data ...any) {
	self.Logger.Warn(msg, data)
}

func (self *GormLogger) Error(c context.Context, msg string, data ...any) {
	self.Logger.Error(msg, data)
}

func (self *GormLogger) Trace(c context.Context, begin time.Time, cb func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := cb()
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")
	sql = strings.TrimSpace(sql)

	if err != nil {
		self.Logger.Error(err.Error(), "sql", sql, "rowsAffected", rowsAffected)
	} else {
		self.Logger.Debug("GORM", "sql", sql, "rowsAffected", rowsAffected)
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
