package validator

import (
	"regexp"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
)

var (
	passwordRegexp = regexp.MustCompile(`[a-zA-Z0-9\~\!\@\#\$\%\^\&\*\(\)\_\-\+\=\{\}\[\]\\\|\:\;\"\'\<\>\,\.\?\/\," "]+`)
)

func init() {
	govalidator.CustomTypeTagMap.Set(
		"password",
		govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
			str, ok := i.(string)
			if !ok {
				return false
			}

			strLength := utf8.RuneCountInString(str)
			if strLength < 6 || strLength > 20 {
				return false
			}

			return passwordRegexp.MatchString(str)
		}),
	)

	govalidator.CustomTypeTagMap.Set(
		"email_or_username",
		govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
			str, ok := i.(string)
			if !ok {
				return false
			}

			if govalidator.IsEmail(str) {
				return true
			}

			strLength := utf8.RuneCountInString(str)
			if strLength < 1 || strLength > 32 {
				return false
			}

			return govalidator.IsUTFLetterNumeric(str)
		}),
	)
}

func ValidateStruct(data interface{}) error {
	if _, err := govalidator.ValidateStruct(data); err != nil {
		return err
	}

	return nil
}
