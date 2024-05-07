package mails

import (
	"fmt"
	"net/smtp"
	"strconv"

	mailConfigurations "github.com/dangduoc08/ecommerce-api/admins/mail_configurations"
	"github.com/dangduoc08/gogo/core"
)

type MailProvider struct{}

func (instance MailProvider) NewProvider() core.Provider {

	return instance
}

func (instance MailProvider) SendMail(
	mailConfiguration *mailConfigurations.MailConfigurationModel,
) func(to string, subject string, messages string) error {
	smtpAuth := smtp.PlainAuth(
		"",
		mailConfiguration.Username,
		mailConfiguration.Password,
		mailConfiguration.Host,
	)

	address := mailConfiguration.Host + ":" + strconv.Itoa(mailConfiguration.Port)

	return func(to, subject, messages string) error {
		headers := make(map[string]string)
		headers["Subject"] = subject
		headers["To"] = to
		headers["MIME-version"] = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

		body := ""
		for k, v := range headers {
			body += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		body += "\r\n" + messages

		err := smtp.SendMail(
			address,
			smtpAuth,
			mailConfiguration.Username,
			[]string{to},
			[]byte(body),
		)

		if err != nil {
			return err
		}

		return nil
	}
}
