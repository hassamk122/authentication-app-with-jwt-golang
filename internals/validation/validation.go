package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(i any) error {
	err := validate.Struct(i)
	if err != nil {
		var errMsgs []string

		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("%s is required", strings.ToLower(e.Field())))
		}

		return fmt.Errorf("%s", strings.Join(errMsgs, ", "))
	}

	return nil
}
