package util

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

// Validate validates a given struct based on its defined tags and returns a map of validation errors if any exist.
func Validate[T any](data T) map[string]string {
	err := validator.New().Struct(data)
	res := make(map[string]string)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			res[err.StructField()] = TranslateTag(err)
		}
	}
	return res
}

// TranslateTag converts a validation error tag into a human-readable error message.
func TranslateTag(fd validator.FieldError) string {
	switch fd.ActualTag() {
	case "required":
		return fmt.Sprintf("%s is required", fd.StructField())
	}
	return "Error validation"
}
