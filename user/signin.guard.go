package user

import (
	"github.com/dangduoc08/ecommerce-api/database"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
)

type SigninGuard struct {
	UserProvider     Provider
	DatabaseProvider database.Provider
	Logger           common.Logger
}

func (signinGuard SigninGuard) CanActivate(ctx gooh.Context) bool {
	body := ctx.Body()

	if !body.Has("data.password") || !body.Has("data.username") {
		return false
	}

	hash := ""
	if username, ok := ctx.Body().Get("data.username").(string); ok {
		user := &User{Username: username}
		resp := signinGuard.DatabaseProvider.DB.
			Where(user).
			First(user)

		// check user exists in db
		if resp.Error != nil {
			signinGuard.Logger.Debug(
				"Error While Query",
				"error", resp.Error.Error(),
				"X-Request-ID", ctx.GetID(),
			)
			return false
		}

		// check user should be active
		if user.Status != ACTIVE {
			return false
		}

		hash = user.Hash
	}

	if password, ok := ctx.Body().Get("data.password").(string); ok {

		// check password
		return signinGuard.UserProvider.CheckPasswordHash(password, hash)
	}

	return false
}
