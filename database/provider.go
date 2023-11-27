package database

import (
	"github.com/dangduoc08/gooh/core"
	"gorm.io/gorm"
)

type Provider struct {
	DB *gorm.DB
}

func (provider Provider) NewProvider() core.Provider {
	return provider
}
