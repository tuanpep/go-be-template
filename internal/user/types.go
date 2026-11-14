package user

import (
	"time"

	"github.com/google/uuid"
)

// User is the domain model. Transport and store types stay at the package edges.
type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateInput is the normalized payload for creating a user.
type CreateInput struct {
	Email string
	Name  string
}

// UpdateInput is a partial update. A nil field is left unchanged.
type UpdateInput struct {
	Email *string
	Name  *string
}
