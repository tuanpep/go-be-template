package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type stubRepo struct {
	email string
	name  string
}

func (s *stubRepo) Get(context.Context, uuid.UUID) (User, error) { return User{}, nil }
func (s *stubRepo) List(context.Context, int32, int32) ([]User, int64, error) {
	return nil, 0, nil
}

func (s *stubRepo) Create(_ context.Context, email, name string) (User, error) {
	s.email = email
	s.name = name
	return User{ID: "id", Email: email, Name: name}, nil
}

func (s *stubRepo) Update(_ context.Context, _ uuid.UUID, email, name *string) (User, error) {
	if email != nil {
		s.email = *email
	}
	if name != nil {
		s.name = *name
	}
	return User{ID: "id", Email: s.email, Name: s.name}, nil
}

func (s *stubRepo) Delete(context.Context, uuid.UUID) error { return nil }

func TestService_Create_NormalizesInput(t *testing.T) {
	t.Parallel()
	repo := &stubRepo{}
	svc := NewService(repo)
	_, err := svc.Create(context.Background(), "  Alice@Example.COM ", "  Bob  ")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if repo.email != "alice@example.com" {
		t.Fatalf("email: got %q", repo.email)
	}
	if repo.name != "Bob" {
		t.Fatalf("name: got %q", repo.name)
	}
}

func TestService_Update_NormalizesEmail(t *testing.T) {
	t.Parallel()
	repo := &stubRepo{}
	svc := NewService(repo)
	email := "  Ada@Example.COM "
	_, err := svc.Update(context.Background(), uuid.New(), &email, nil)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if repo.email != "ada@example.com" {
		t.Fatalf("email: got %q", repo.email)
	}
}

func TestService_Update_RequiresAField(t *testing.T) {
	t.Parallel()
	_, err := NewService(&stubRepo{}).Update(context.Background(), uuid.New(), nil, nil)
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_Create_InvalidInput(t *testing.T) {
	t.Parallel()
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), "   ", "x")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}
