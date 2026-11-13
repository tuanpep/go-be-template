package user

import (
	"fmt"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/tuanbt/go-be-template/internal/transport/openapi"
)

// ToAPI maps a domain user to the OpenAPI response model.
func ToAPI(u User) (openapi.User, error) {
	id, err := uuid.Parse(u.ID)
	if err != nil {
		return openapi.User{}, fmt.Errorf("parse id: %w", err)
	}
	return openapi.User{
		Id:        openapi_types.UUID(id),
		Email:     openapi_types.Email(u.Email),
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

// ParseID converts an OpenAPI UUID path param to a domain id.
func ParseID(id openapi_types.UUID) (uuid.UUID, error) {
	return uuid.Parse(id.String())
}
