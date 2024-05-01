package permissions

import (
	"database/sql/driver"
	"strings"
)

type PermissionModel []string

func (instance *PermissionModel) Scan(src any) error {
	*instance = strings.Split(src.(string), ",")
	return nil
}

func (instance PermissionModel) Value() (driver.Value, error) {
	if len(instance) == 0 {
		return nil, nil
	}
	return strings.Join(instance, ","), nil
}
