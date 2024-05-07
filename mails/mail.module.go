package mails

import (
	"github.com/dangduoc08/gogo/core"
)

var MailModule = func() *core.Module {
	module := core.ModuleBuilder().
		Providers(
			MailProvider{},
		).
		Build()

	return module
}()
