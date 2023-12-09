package users

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/db"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
)

type SessionsGuard struct {
	UserProvider UserProvider
	DBProvider   db.DBProvider
	Logger       common.Logger
}

func (self SessionsGuard) CanActivate(ctx gooh.Context) bool {
	body := ctx.Body()

	if !body.Has("data.password") || !body.Has("data.username") {
		return false
	}

	hash := ""
	var id uint

	if username, ok := body.Get("data.username").(string); ok {
		user := &User{Username: username}
		resp := self.DBProvider.DB.
			Where(user).
			First(user)

		// check user exists in db
		if resp.Error != nil {
			self.Logger.Debug(
				"Error While Query",
				"error", resp.Error.Error(),
				"X-Request-ID", ctx.GetID(),
			)
			return false
		}

		// check user should be active
		if user.Status != UserStatus(ACTIVE) {
			panic(exception.UnauthorizedException(fmt.Sprintf("Field: user.status, Error: %s", user.Status)))
		}

		hash = user.Hash
		id = user.ID
	}

	if password, ok := body.Get("data.password").(string); ok {
		body.Set("data.ID", float64(id))

		// check password
		return self.UserProvider.CheckPasswordHash(password, hash)
	}

	return false
}
