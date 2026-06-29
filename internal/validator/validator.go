package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseValidationErrors(err error) string {
	if err == nil {
		return ""
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return "Invalid request."
	}

	var messages []string

	for _, fieldErr := range validationErrors {
		field := fieldErr.Field()

		var message string

		switch fieldErr.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required.", field)

		case "email":
			message = fmt.Sprintf("%s must be a valid email address.", field)

		case "min":
			message = fmt.Sprintf("%s must be at least %s characters long.", field, fieldErr.Param())

		case "max":
			message = fmt.Sprintf("%s must not exceed %s characters.", field, fieldErr.Param())

		default:
			message = fmt.Sprintf("%s is invalid.", field)
		}

		messages = append(messages, "• "+message)
	}

	return strings.Join(messages, "\n")
}
