package user

import (
	"github.com/dangduoc08/ecommerce-api/database"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/exception"
)

type CreateGuard struct {
	UserProvider     Provider
	DatabaseProvider database.Provider
}

func (createGuard CreateGuard) CanActivate(ctx gooh.Context) bool {
	body := ctx.Body()

	data := []map[string]string{
		{
			"email": body.Get("data.email").(string),
		},
		{
			"username": body.Get("data.username").(string),
		},
	}
	err := createGuard.UserProvider.CheckDuplicateUser(data)
	if err != nil {
		errMsgs := []string{
			"Field: user.email or user.username",
			"Error: duplicate data",
		}
		panic(exception.ConflictException(errMsgs))
	}

	return true
}
