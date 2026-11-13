package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/tuanbt/go-be-template/internal/store/db"
)

// Repository persists users (interface at the consumer — the service).
type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (User, error)
	List(ctx context.Context, limit, offset int32) ([]User, int64, error)
	Create(ctx context.Context, email, name string) (User, error)
}

// SQLRepository implements Repository using sqlc-generated queries.
type SQLRepository struct {
	q *db.Queries
}

func NewSQLRepository(q *db.Queries) *SQLRepository {
	return &SQLRepository{q: q}
}

func (r *SQLRepository) Get(ctx context.Context, id uuid.UUID) (User, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id.String()); err != nil {
		return User{}, fmt.Errorf("scan id: %w", err)
	}
	row, err := r.q.GetUser(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("get user: %w", err)
	}
	return fromDB(row)
}

func (r *SQLRepository) List(ctx context.Context, limit, offset int32) ([]User, int64, error) {
	rows, err := r.q.ListUsers(ctx, db.ListUsersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	total, err := r.q.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	out := make([]User, 0, len(rows))
	for _, row := range rows {
		u, err := fromDB(row)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, u)
	}
	return out, total, nil
}

func (r *SQLRepository) Create(ctx context.Context, email, name string) (User, error) {
	row, err := r.q.CreateUser(ctx, db.CreateUserParams{Email: email, Name: name})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
			return User{}, ErrDuplicateEmail
		}
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return fromDB(row)
}

func fromDB(row db.User) (User, error) {
	if !row.ID.Valid || !row.CreatedAt.Valid || !row.UpdatedAt.Valid {
		return User{}, fmt.Errorf("incomplete user row")
	}
	id, err := uuid.FromBytes(row.ID.Bytes[:])
	if err != nil {
		return User{}, fmt.Errorf("user id: %w", err)
	}
	return User{
		ID:        id.String(),
		Email:     row.Email,
		Name:      row.Name,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
