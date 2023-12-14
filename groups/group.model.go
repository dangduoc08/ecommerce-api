package groups

import (
	"database/sql/driver"
	"strings"
	"time"
)

type Permission []string

func (permission *Permission) Scan(src any) error {
	*permission = strings.Split(src.(string), ",")
	return nil
}

func (permission Permission) Value() (driver.Value, error) {
	if len(permission) == 0 {
		return nil, nil
	}
	return strings.Join(permission, ","), nil
}

type Group struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name"`
	Permissions Permission `json:"permissions" gorm:"type:TEXT"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime:true"`
}
