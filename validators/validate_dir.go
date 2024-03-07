package validators

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

func ValidateDir(cb func(validator.FieldError)) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value != "" {
			if parsedURL, err := url.Parse(value); err != nil {
				cb(FieldError{
					fieldLevel: fl,
					val:        value,
				})
				return true
			} else {
				currentDir, _ := os.Getwd()
				publicPath := filepath.Join(currentDir, "publics", parsedURL.Path)

				if parsedURL.Host == "" {
					if _, err := os.Stat(publicPath); err != nil {
						cb(FieldError{
							fieldLevel: fl,
							val:        value,
						})
						return true
					}
				}
			}
		}

		return true
	}
}
