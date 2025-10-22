package interfaces

import (
	"context"

	"github.com/ibnuzaman/ewallet-ums/internal/models"
)

// IUserRepository defines the interface for user repository operations.
type IUserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*models.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// GetByPhone retrieves a user by phone
	GetByPhone(ctx context.Context, phone string) (*models.User, error)

	// Update updates a user
	Update(ctx context.Context, user *models.User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id int64) error

	// List retrieves users based on filters
	List(ctx context.Context, filter models.UserFilter) ([]*models.User, error)

	// Count counts users based on filters
	Count(ctx context.Context, filter models.UserFilter) (int64, error)
}
