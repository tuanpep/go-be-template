package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tuanbt/go-be-template/internal/store/db"
)

// Repository persists users (interface at the consumer — the service).
type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (User, error)
	List(ctx context.Context, limit, offset int32) ([]User, int64, error)
	Create(ctx context.Context, email, name string) (User, error)
	Update(ctx context.Context, id uuid.UUID, email, name *string) (User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// SQLRepository implements Repository using sqlc-generated queries.
type SQLRepository struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewSQLRepository(pool *pgxpool.Pool, q *db.Queries) *SQLRepository {
	return &SQLRepository{pool: pool, q: q}
}

func (r *SQLRepository) Get(ctx context.Context, id uuid.UUID) (User, error) {
	pgID, err := pgUUID(id)
	if err != nil {
		return User{}, err
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
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return nil, 0, fmt.Errorf("begin list: %w", err)
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	rows, err := q.ListUsers(ctx, db.ListUsersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	total, err := q.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, 0, fmt.Errorf("commit list: %w", err)
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
		if duplicateEmail(err) {
			return User{}, ErrDuplicateEmail
		}
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return fromDB(row)
}

func (r *SQLRepository) Update(ctx context.Context, id uuid.UUID, email, name *string) (User, error) {
	pgID, err := pgUUID(id)
	if err != nil {
		return User{}, err
	}
	row, err := r.q.UpdateUser(ctx, db.UpdateUserParams{ID: pgID, Email: email, Name: name})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}
		if duplicateEmail(err) {
			return User{}, ErrDuplicateEmail
		}
		return User{}, fmt.Errorf("update user: %w", err)
	}
	return fromDB(row)
}

func (r *SQLRepository) Delete(ctx context.Context, id uuid.UUID) error {
	pgID, err := pgUUID(id)
	if err != nil {
		return err
	}
	n, err := r.q.DeleteUser(ctx, pgID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func pgUUID(id uuid.UUID) (pgtype.UUID, error) {
	var pgID pgtype.UUID
	if err := pgID.Scan(id.String()); err != nil {
		return pgtype.UUID{}, fmt.Errorf("scan id: %w", err)
	}
	return pgID, nil
}

func duplicateEmail(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key"
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
