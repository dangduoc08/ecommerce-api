package mailConfigurations

import (
	"github.com/dangduoc08/gogo/core"
)

var MailConfigurationModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			MailConfigurationController{},
		).
		Providers(
			MailConfigurationProvider{},
		).
		Build()

	module.
		Prefix("mail_configurations")

	return module
}()
