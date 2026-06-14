package service

import (
	"context"
	"errors"
	"time"

	"github.com/example/user-age-api/db/sqlc"
	"github.com/example/user-age-api/internal/models"
	"github.com/example/user-age-api/internal/repository"
	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found")

// UserService contains business logic for user operations.
type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error)
	GetUser(ctx context.Context, id int32) (models.UserWithAgeResponse, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context, page, pageSize int32) (models.PaginatedUsersResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse(models.DateLayout, req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Create(ctx, sqlc.CreateUserParams{
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format(models.DateLayout),
	}, nil
}

func (s *userService) GetUser(ctx context.Context, id int32) (models.UserWithAgeResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserWithAgeResponse{}, ErrUserNotFound
		}
		return models.UserWithAgeResponse{}, err
	}

	return models.UserWithAgeResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format(models.DateLayout),
		Age:  models.CalculateAge(user.Dob),
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	dob, err := time.Parse(models.DateLayout, req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.Update(ctx, sqlc.UpdateUserParams{
		ID:   id,
		Name: req.Name,
		Dob:  dob,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserResponse{}, ErrUserNotFound
		}
		return models.UserResponse{}, err
	}

	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.Format(models.DateLayout),
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	rows, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *userService) ListUsers(ctx context.Context, page, pageSize int32) (models.PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	users, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		return models.PaginatedUsersResponse{}, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return models.PaginatedUsersResponse{}, err
	}

	result := make([]models.UserWithAgeResponse, 0, len(users))
	for _, u := range users {
		result = append(result, models.UserWithAgeResponse{
			ID:   u.ID,
			Name: u.Name,
			Dob:  u.Dob.Format(models.DateLayout),
			Age:  models.CalculateAge(u.Dob),
		})
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}

	return models.PaginatedUsersResponse{
		Data:       result,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	}, nil
}
