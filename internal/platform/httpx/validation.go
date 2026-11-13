package httpx

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ValidationDetail returns a client-safe 400 message for validator failures.
func ValidationDetail(err error) string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) || len(verrs) == 0 {
		return "invalid request"
	}
	e := verrs[0]
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldLabel(e.Field()))
	case "email":
		return "invalid email address"
	case "min", "max":
		return fmt.Sprintf("%s is out of range", fieldLabel(e.Field()))
	default:
		return "invalid request"
	}
}

func fieldLabel(field string) string {
	switch field {
	case "Email":
		return "email"
	case "Name":
		return "name"
	default:
		return "field"
	}
}
