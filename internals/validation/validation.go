package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(i any) error {

	err := validate.Struct(i)

	if err != nil {
		var errMsgs string

		for _, e := range err.(validator.ValidationErrors) {
			switch e.Field() {
			case "Name":
				errMsgs = "Name must be between 2 and 50 characters"
				return fmt.Errorf("%s", errMsgs)
			case "Email":
				errMsgs = "Must be a valid email address"
				return fmt.Errorf("%s", errMsgs)
			case "Password":
				errMsgs = "Password must be at least 8 characters"
				return fmt.Errorf("%s", errMsgs)
			}
		}

	}

	return nil
}
