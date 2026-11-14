package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Service implements user use cases.
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (User, error) {
	u, err := s.repo.Get(ctx, id)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *Service) List(ctx context.Context, limit, offset int32) ([]User, int64, error) {
	users, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list: %w", err)
	}
	return users, total, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, email, name *string) (User, error) {
	if email == nil && name == nil {
		return User{}, ErrInvalidInput
	}
	if email != nil {
		v := strings.ToLower(strings.TrimSpace(*email))
		if v == "" {
			return User{}, ErrInvalidInput
		}
		email = &v
	}
	if name != nil {
		v := strings.TrimSpace(*name)
		if v == "" {
			return User{}, ErrInvalidInput
		}
		name = &v
	}
	u, err := s.repo.Update(ctx, id, email, name)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) Create(ctx context.Context, email, name string) (User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)
	if email == "" || name == "" {
		return User{}, ErrInvalidInput
	}
	u, err := s.repo.Create(ctx, email, name)
	if err != nil {
		return User{}, err
	}
	return u, nil
}
