package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ibnuzaman/ewallet-ums/helpers"
	"github.com/ibnuzaman/ewallet-ums/internal/models"
)

// UserRepository implements IUserRepository.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user.
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, phone, full_name, password_hash, is_active, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowxContext(
		ctx,
		query,
		user.Email,
		user.Phone,
		user.FullName,
		user.PasswordHash,
		user.IsActive,
		user.IsVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		helpers.Logger.Errorf("Failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	helpers.Logger.Infof("User created successfully with ID: %d", user.ID)
	return nil
}

// GetByID retrieves a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, email, phone, full_name, password_hash, is_active, is_verified,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		helpers.Logger.Errorf("Failed to get user by ID %d: %v", id, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, phone, full_name, password_hash, is_active, is_verified,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		helpers.Logger.Errorf("Failed to get user by email %s: %v", email, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByPhone retrieves a user by phone.
func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `
		SELECT id, email, phone, full_name, password_hash, is_active, is_verified,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE phone = $1 AND deleted_at IS NULL
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		helpers.Logger.Errorf("Failed to get user by phone %s: %v", phone, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// Update updates a user.
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, phone = $2, full_name = $3, is_active = $4, is_verified = $5, updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Email,
		user.Phone,
		user.FullName,
		user.IsActive,
		user.IsVerified,
		time.Now(),
		user.ID,
	)
	if err != nil {
		helpers.Logger.Errorf("Failed to update user %d: %v", user.ID, err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	helpers.Logger.Infof("User %d updated successfully", user.ID)
	return nil
}

// Delete soft deletes a user.
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE users
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		helpers.Logger.Errorf("Failed to delete user %d: %v", id, err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	helpers.Logger.Infof("User %d deleted successfully", id)
	return nil
}

// List retrieves users based on filters.
func (r *UserRepository) List(ctx context.Context, filter models.UserFilter) ([]*models.User, error) {
	query := `
		SELECT id, email, phone, full_name, password_hash, is_active, is_verified,
		       created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
	`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, filter.Email)
		argIndex++
	}

	if filter.Phone != "" {
		conditions = append(conditions, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, filter.Phone)
		argIndex++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.IsVerified != nil {
		conditions = append(conditions, fmt.Sprintf("is_verified = $%d", argIndex))
		args = append(args, *filter.IsVerified)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	var users []*models.User
	err := r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		helpers.Logger.Errorf("Failed to list users: %v", err)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// Count counts users based on filters.
func (r *UserRepository) Count(ctx context.Context, filter models.UserFilter) (int64, error) {
	query := "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL"

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Email != "" {
		conditions = append(conditions, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, filter.Email)
		argIndex++
	}

	if filter.Phone != "" {
		conditions = append(conditions, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, filter.Phone)
		argIndex++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if filter.IsVerified != nil {
		conditions = append(conditions, fmt.Sprintf("is_verified = $%d", argIndex))
		args = append(args, *filter.IsVerified)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var count int64
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		helpers.Logger.Errorf("Failed to count users: %v", err)
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}
