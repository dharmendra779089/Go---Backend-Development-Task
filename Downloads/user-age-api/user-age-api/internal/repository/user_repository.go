package repository

import (
	"context"

	"github.com/example/user-age-api/db/sqlc"
)

// UserRepository defines data-access operations for users.
type UserRepository interface {
	Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	GetByID(ctx context.Context, id int32) (sqlc.User, error)
	Update(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)
	Delete(ctx context.Context, id int32) (int64, error)
	List(ctx context.Context, limit, offset int32) ([]sqlc.User, error)
	Count(ctx context.Context) (int64, error)
}

type userRepository struct {
	queries *sqlc.Queries
}

// NewUserRepository creates a new repository backed by sqlc generated queries.
func NewUserRepository(queries *sqlc.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, arg)
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (sqlc.User, error) {
	return r.queries.GetUser(ctx, id)
}

func (r *userRepository) Update(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	return r.queries.UpdateUser(ctx, arg)
}

func (r *userRepository) Delete(ctx context.Context, id int32) (int64, error) {
	return r.queries.DeleteUser(ctx, id)
}

func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]sqlc.User, error) {
	return r.queries.ListUsers(ctx, sqlc.ListUsersParams{Limit: limit, Offset: offset})
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}
