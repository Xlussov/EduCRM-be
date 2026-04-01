package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{
		validator: v,
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func ParseError(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return map[string]string{"error": err.Error()}
	}

	errs := make(map[string]string)
	for _, f := range validationErrors {
		field := f.Field()

		switch f.Tag() {
		case "required":
			errs[field] = "This field is required"
		case "min":
			errs[field] = fmt.Sprintf("Must be at least %s characters long", f.Param())
		case "max":
			errs[field] = fmt.Sprintf("Must be a maximum of %s characters long", f.Param())
		case "email":
			errs[field] = "Invalid email format"
		case "e164":
			errs[field] = "Invalid phone number format (use +1234567890)"
		case "datetime":
			errs[field] = fmt.Sprintf("Invalid date format, expected %s", f.Param())
		case "oneof":
			errs[field] = fmt.Sprintf("Must be one of: %s", strings.ReplaceAll(f.Param(), " ", ", "))
		case "jwt":
			errs[field] = "Invalid JWT token"
		case "unique":
			errs[field] = "Value must be unique"
		case "uuid":
			errs[field] = "Invalid UUID format"
		default:
			errs[field] = "Invalid value"
		}
	}

	return errs
}
