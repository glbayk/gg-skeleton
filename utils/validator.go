package utils

import (
	"github.com/go-playground/validator/v10"
)

func Validate[T any](payload T) (errors []string) {
	err := validator.New().Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Error())
		}
	}

	return errors
}
