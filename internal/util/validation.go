package util

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// Validate validates a given struct based on its defined tags and returns a map of validation errors if any exist.
func Validate[T any](data T) map[string]string {
	logrus.Infof("Validating struct: %+v", data)
	err := validator.New().Struct(data)
	res := make(map[string]string)
	if err != nil {
		logrus.Warnf("Validation errors found: %v", err)
		for _, err := range err.(validator.ValidationErrors) {
			res[err.Field()] = TranslateTag(err)
		}
	}
	logrus.Infof("Validation result: %+v", res)
	return res
}

// TranslateTag converts a validation error tag into a human-readable error message.
func TranslateTag(fd validator.FieldError) string {
	switch fd.ActualTag() {
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fd.StructField(), fd.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fd.StructField(), fd.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", fd.StructField(), fd.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", fd.StructField(), fd.Param())
	case "numeric":
		return fmt.Sprintf("%s must be a number", fd.StructField())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fd.StructField())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fd.StructField())
	case "len":
		return fmt.Sprintf("%s must be %s characters long", fd.StructField(), fd.Param())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fd.StructField(), fd.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fd.StructField(), fd.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", fd.StructField(), fd.Param())
	case "required":
		return fmt.Sprintf("%s is required", fd.StructField())
	default:
		return fmt.Sprintf("Invalid value for %s", fd.StructField())
	}
}
