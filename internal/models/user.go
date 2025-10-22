package models

import (
	"database/sql"
	"time"
)

// User represents a user in the system.
type User struct {
	PasswordHash string       `db:"password_hash" json:"-"`
	Email        string       `db:"email" json:"email"`
	Phone        string       `db:"phone" json:"phone"`
	FullName     string       `db:"full_name" json:"full_name"`
	CreatedAt    time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt    sql.NullTime `db:"deleted_at" json:"deleted_at,omitempty"`
	ID           int64        `db:"id" json:"id"`
	IsActive     bool         `db:"is_active" json:"is_active"`
	IsVerified   bool         `db:"is_verified" json:"is_verified"`
}

// CreateUserRequest represents the request to create a user.
type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest represents the request to update a user.
type UpdateUserRequest struct {
	FullName *string `json:"full_name,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}

// UserFilter represents filters for querying users.
//
//nolint:govet // fieldalignment: reordering would hurt readability
type UserFilter struct {
	Email      string
	Phone      string
	IsActive   *bool
	IsVerified *bool
	Offset     int
	Limit      int
}
