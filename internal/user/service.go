package user

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

// Service is the feature's public use-case surface.
// Callers depend on this interface. The implementation stays in this package.
type Service interface {
	Get(ctx context.Context, id uuid.UUID) (User, error)
	List(ctx context.Context, limit, offset int32) ([]User, int64, error)
	Create(ctx context.Context, in CreateInput) (User, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateInput) (User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// Repository is the persistence port consumed by the service.
type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (User, error)
	List(ctx context.Context, limit, offset int32) ([]User, int64, error)
	Create(ctx context.Context, in CreateInput) (User, error)
	Update(ctx context.Context, id uuid.UUID, in UpdateInput) (User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

// NewService returns the user use cases.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, id uuid.UUID) (User, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) List(ctx context.Context, limit, offset int32) ([]User, int64, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (User, error) {
	if in.Email == nil && in.Name == nil {
		return User{}, ErrInvalidInput
	}
	if in.Email != nil {
		v := strings.ToLower(strings.TrimSpace(*in.Email))
		if v == "" {
			return User{}, ErrInvalidInput
		}
		in.Email = &v
	}
	if in.Name != nil {
		v := strings.TrimSpace(*in.Name)
		if v == "" {
			return User{}, ErrInvalidInput
		}
		in.Name = &v
	}
	return s.repo.Update(ctx, id, in)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) Create(ctx context.Context, in CreateInput) (User, error) {
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.Name = strings.TrimSpace(in.Name)
	if in.Email == "" || in.Name == "" {
		return User{}, ErrInvalidInput
	}
	return s.repo.Create(ctx, in)
}
