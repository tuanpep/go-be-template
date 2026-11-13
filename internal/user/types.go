package user

import "time"

// User is the domain model (transport/store mapping stays at package edges).
type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
