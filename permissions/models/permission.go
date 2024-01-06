package models

import (
	"database/sql/driver"
	"strings"
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
