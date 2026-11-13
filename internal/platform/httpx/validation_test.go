package httpx

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestValidationDetail(t *testing.T) {
	t.Parallel()
	v := validator.New()
	err := v.Struct(struct {
		Email string `validate:"required,email"`
		Name  string `validate:"required,min=1"`
	}{Email: "bad", Name: ""})
	if err == nil {
		t.Fatal("expected validation error")
	}
	msg := ValidationDetail(err)
	if msg == "" || msg == "invalid request body" {
		t.Fatalf("unexpected detail: %q", msg)
	}
}
