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

func (s *stubRepo) Create(_ context.Context, in CreateInput) (User, error) {
	s.email = in.Email
	s.name = in.Name
	return User{Email: in.Email, Name: in.Name}, nil
}

func (s *stubRepo) Update(_ context.Context, _ uuid.UUID, in UpdateInput) (User, error) {
	if in.Email != nil {
		s.email = *in.Email
	}
	if in.Name != nil {
		s.name = *in.Name
	}
	return User{Email: s.email, Name: s.name}, nil
}

func (s *stubRepo) Delete(context.Context, uuid.UUID) error { return nil }

func TestService_Create_NormalizesInput(t *testing.T) {
	t.Parallel()
	repo := &stubRepo{}
	svc := NewService(repo)
	_, err := svc.Create(context.Background(), CreateInput{Email: "  Alice@Example.COM ", Name: "  Bob  "})
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
	_, err := svc.Update(context.Background(), uuid.New(), UpdateInput{Email: &email})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if repo.email != "ada@example.com" {
		t.Fatalf("email: got %q", repo.email)
	}
}

func TestService_Update_RequiresAField(t *testing.T) {
	t.Parallel()
	_, err := NewService(&stubRepo{}).Update(context.Background(), uuid.New(), UpdateInput{})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_Create_InvalidInput(t *testing.T) {
	t.Parallel()
	svc := NewService(&stubRepo{})
	_, err := svc.Create(context.Background(), CreateInput{Email: "   ", Name: "x"})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}
