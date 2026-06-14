package models

import "time"

// CreateUserRequest is the payload for creating a user.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	Dob  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest is the payload for updating a user.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	Dob  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UserResponse is returned for create/update (no age field).
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
}

// UserWithAgeResponse is returned for get/list (includes age).
type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
	Age  int    `json:"age"`
}

// PaginatedUsersResponse wraps a list of users with pagination metadata.
type PaginatedUsersResponse struct {
	Data       []UserWithAgeResponse `json:"data"`
	Page       int32                 `json:"page"`
	PageSize   int32                 `json:"page_size"`
	TotalCount int64                 `json:"total_count"`
	TotalPages int64                 `json:"total_pages"`
}

const DateLayout = "2006-01-02"

// CalculateAge returns the age in full years for a given dob, as of now (UTC).
func CalculateAge(dob time.Time) int {
	now := time.Now().UTC()
	years := now.Year() - dob.Year()

	// If the birthday hasn't occurred yet this year, subtract one year.
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}

	if years < 0 {
		return 0
	}

	return years
}
