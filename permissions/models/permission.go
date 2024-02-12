package models

import (
	"database/sql/driver"
	"strings"
)

type Permissions []string

func (permissions *Permissions) Scan(src any) error {
	*permissions = strings.Split(src.(string), ",")
	return nil
}

func (permissions Permissions) Value() (driver.Value, error) {
	if len(permissions) == 0 {
		return nil, nil
	}
	return strings.Join(permissions, ","), nil
}
