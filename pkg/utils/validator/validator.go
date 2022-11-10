package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)
	if err != nil {
		validationErr := err.(validator.ValidationErrors)
		for _, each := range validationErr {
			switch each.Tag() {
			case "required":
				msg := fmt.Sprintf("%s is required", each.Field())
				return errors.New(msg)
			case "len":
				msg := fmt.Sprintf("%s must be %s characters long", each.Field(), each.Param())
				return errors.New(msg)
			case "gte":
				msg := fmt.Sprintf("%s must be greater than or equal to %s", each.Field(), each.Param())
				return errors.New(msg)
			default:
				msg := fmt.Sprintf("Invalid field %s", each.Field())
				return errors.New(msg)
			}
		}
	}

	return nil
}
