package auths

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/dangduoc08/ecommerce-api/admins/groups"
	mailConfigurations "github.com/dangduoc08/ecommerce-api/admins/mail_configurations"
	"github.com/dangduoc08/ecommerce-api/admins/stores"
	"github.com/dangduoc08/ecommerce-api/admins/users"
	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/mails"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/modules/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthProvider struct {
	config.ConfigService
	common.Logger

	TemplatePath          string
	MailProvider          mails.MailProvider
	JWTResetPasswordKey   string
	JWTResetPasswordExpIn int
}

func (instance AuthProvider) NewProvider() core.Provider {
	currentDir, _ := os.Getwd()
	instance.TemplatePath = filepath.Join(currentDir, constants.TEMPLATE_DIR, "reset_password_email.html")

	instance.JWTResetPasswordKey = instance.Get("JWT_RESET_PASSWORD_KEY").(string)
	instance.JWTResetPasswordExpIn = instance.Get("JWT_RESET_PASSWORD_EXP_IN").(int)

	return instance
}

func (instance AuthProvider) GetUserPermissions(groups []*groups.GroupModel) []string {
	permissions := []string{}

	for _, group := range groups {
		permissions = append(permissions, group.Permissions...)
	}

	return utils.ArrToUnique(permissions)
}

func (instance AuthProvider) HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(hashedPasswordBytes), nil
}

func (instance AuthProvider) CheckHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (instance AuthProvider) SignToken(claims jwt.MapClaims, key string, expIn int) (string, error) {
	claims["exp"] = time.Now().Unix() + int64(expIn)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

func (instance AuthProvider) SendResetPasswordEmail(
	domain string,
	mailConfiguration *mailConfigurations.MailConfigurationModel,
	store *stores.StoreModel,
	user *users.UserModel,
) error {
	tmpl, err := template.ParseFiles(instance.TemplatePath)
	if err != nil {
		return err
	}

	resetURL, err := url.Parse(domain)
	if err != nil {
		return err
	}

	resetPasswordToken, err := instance.SignToken(
		jwt.MapClaims{
			"user_id":  user.ID,
			"store_id": user.StoreID,
		},
		instance.JWTResetPasswordKey,
		instance.JWTResetPasswordExpIn,
	)
	if err != nil {
		return err
	}

	query := resetURL.Query()
	query.Set("token", resetPasswordToken)

	resetURL = resetURL.JoinPath("/password-reset")
	resetURL.RawQuery = query.Encode()

	msgPayload := map[string]string{
		"Title":        constants.RESET_PASSWORD_SUBJECT,
		"ResetURL":     resetURL.String(),
		"FirstName":    user.FirstName,
		"LastName":     user.LastName,
		"ContactEmail": store.Email,
	}

	var msgBuf bytes.Buffer
	if err := tmpl.Execute(&msgBuf, msgPayload); err != nil {
		return err
	}

	if err = instance.MailProvider.SendMail(mailConfiguration)(
		user.Email,
		constants.RESET_PASSWORD_SUBJECT,
		msgBuf.String(),
	); err != nil {
		return err
	}

	return nil
}
