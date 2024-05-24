package utils

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var uuidRegex *regexp.Regexp

func init() {
	uuidRegex, _ = regexp.Compile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
}

func RegisterCustomValidator(validate *validator.Validate) {
	err := validate.RegisterValidation("uuid", Uuid)
	if err != nil {
		panic(err)
	}
}

func ParseValidatorErr(err error) []string {
	errMessages := []string{}

	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		field := e.Field()
		if len(field) == 0 {
			field = "param"
		}
		errMessages = append(errMessages, fmt.Sprintf("%s failed on %s", field, e.Tag()))
	}

	return errMessages
}

func Uuid(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return uuidRegex.MatchString(value)
}
