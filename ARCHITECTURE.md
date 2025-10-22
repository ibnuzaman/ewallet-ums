# Clean Architecture - Database Implementation

## Overview

Implementasi database PostgreSQL dengan clean architecture pattern untuk ewallet-ums.

## Architecture Layers

```
┌─────────────────────────────────────────┐
│           HTTP Handler (API)            │
│         internal/api/*.go               │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│          Service Layer                  │
│       internal/services/*.go            │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│        Repository Layer                 │
│      internal/repository/*.go           │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Database Layer                  │
│        database/postgres.go             │
└─────────────────────────────────────────┘
```

## Components

### 1. Database Layer (`database/postgres.go`)

**Responsibility**: Mengelola koneksi database dan connection pool

**Key Features**:
- Singleton pattern untuk koneksi database
- Connection pooling dengan konfigurasi optimal
- Health check function
- Graceful connection management

**Functions**:
```go
InitPostgres() (*sqlx.DB, error)    // Initialize connection
GetPostgresDB() *sqlx.DB             // Get connection instance
ClosePostgres() error                // Close connection
HealthCheck(ctx) error               // Check connection health
```

**Configuration**:
```go
type PostgresConfig struct {
    Host            string
    Port            string
    User            string
    Password        string
    DBName          string
    SSLMode         string
    MaxOpenConns    int           // 25
    MaxIdleConns    int           // 25
    ConnMaxLifetime time.Duration // 5 minutes
    ConnMaxIdleTime time.Duration // 5 minutes
}
```

### 2. Models Layer (`internal/models/`)

**Responsibility**: Define data structures

**Example**:
```go
type User struct {
    ID           int64        `db:"id" json:"id"`
    Email        string       `db:"email" json:"email"`
    Phone        string       `db:"phone" json:"phone"`
    FullName     string       `db:"full_name" json:"full_name"`
    PasswordHash string       `db:"password_hash" json:"-"`
    IsActive     bool         `db:"is_active" json:"is_active"`
    IsVerified   bool         `db:"is_verified" json:"is_verified"`
    CreatedAt    time.Time    `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time    `db:"updated_at" json:"updated_at"`
    DeletedAt    sql.NullTime `db:"deleted_at" json:"deleted_at,omitempty"`
}
```

### 3. Repository Interface (`internal/interfaces/`)

**Responsibility**: Define contract for data operations

**Example**:
```go
type IUserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id int64) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id int64) error
    List(ctx context.Context, filter models.UserFilter) ([]*models.User, error)
    Count(ctx context.Context, filter models.UserFilter) (int64, error)
}
```

### 4. Repository Implementation (`internal/repository/`)

**Responsibility**: Implement data access operations

**Key Features**:
- Context-aware queries (with timeout)
- Proper error handling and logging
- Soft delete support
- Pagination support
- Filter support

**Example**:
```go
type UserRepository struct {
    db *sqlx.DB
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (email, phone, full_name, password_hash)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `
    // Implementation...
}
```

### 5. Service Layer (`internal/services/`)

**Responsibility**: Business logic dan orchestration

**Example**:
```go
type UserService struct {
    userRepo interfaces.IUserRepository
}

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
    // Business logic
    // Validation
    // Call repository
}
```

### 6. API Layer (`internal/api/`)

**Responsibility**: HTTP handlers

**Example**:
```go
type UserAPI struct {
    userService interfaces.IUserService
}

func (api *UserAPI) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request
    // Call service
    // Send response
}
```

## Database Operations

### CRUD Operations

#### Create
```go
user := &models.User{
    Email:        "user@example.com",
    Phone:        "+6281234567890",
    FullName:     "John Doe",
    PasswordHash: hashedPassword,
    IsActive:     true,
    IsVerified:   false,
}

err := userRepo.Create(ctx, user)
// user.ID, user.CreatedAt, user.UpdatedAt will be populated
```

#### Read
```go
// By ID
user, err := userRepo.GetByID(ctx, 1)

// By Email
user, err := userRepo.GetByEmail(ctx, "user@example.com")

// By Phone
user, err := userRepo.GetByPhone(ctx, "+6281234567890")

// List with filters
filter := models.UserFilter{
    IsActive:   boolPtr(true),
    Limit:      10,
    Offset:     0,
}
users, err := userRepo.List(ctx, filter)

// Count
count, err := userRepo.Count(ctx, filter)
```

#### Update
```go
user.FullName = "Jane Doe"
user.IsVerified = true

err := userRepo.Update(ctx, user)
```

#### Delete (Soft Delete)
```go
err := userRepo.Delete(ctx, userID)
```

## Best Practices

### 1. Always Use Context

```go
// Good
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := repo.GetByID(ctx, id)
```

### 2. Handle Errors Properly

```go
user, err := repo.GetByID(ctx, id)
if err != nil {
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    logger.Errorf("Failed to get user: %v", err)
    return nil, fmt.Errorf("failed to get user: %w", err)
}
```

### 3. Use Prepared Statements

sqlx automatically uses prepared statements for repeated queries.

### 4. Transaction Support

```go
func (r *UserRepository) CreateUserWithProfile(ctx context.Context, user *models.User, profile *models.Profile) error {
    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Insert user
    err = tx.QueryRowxContext(ctx, insertUserQuery, ...).Scan(&user.ID)
    if err != nil {
        return err
    }

    // Insert profile
    profile.UserID = user.ID
    _, err = tx.ExecContext(ctx, insertProfileQuery, ...)
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

### 5. Connection Pool Management

```go
// Check pool stats
stats := db.Stats()
logger.Infof("Open connections: %d", stats.OpenConnections)
logger.Infof("In use: %d", stats.InUse)
logger.Infof("Idle: %d", stats.Idle)
```

## Dependency Injection

### In main.go

```go
func main() {
    // Initialize database
    db, err := database.InitPostgres()
    if err != nil {
        log.Fatal(err)
    }
    defer database.ClosePostgres()

    // Initialize repositories
    userRepo := repository.NewUserRepository(db)

    // Initialize services
    userService := services.NewUserService(userRepo)

    // Initialize APIs
    userAPI := api.NewUserAPI(userService)

    // Setup routes
    r.Post("/users", userAPI.CreateUserHandler)
}
```

### In cmd/http.go

```go
func dependencyInject() Dependency {
    db := database.GetPostgresDB()

    // Repositories
    userRepo := repository.NewUserRepository(db)

    // Services
    userService := services.NewUserService(userRepo)

    // APIs
    userAPI := api.NewUserAPI(userService)

    return Dependency{
        UserAPI: userAPI,
    }
}
```

## Testing

### Unit Tests (Repository)

```go
func TestUserRepository_Create(t *testing.T) {
    // Use testcontainers or mock database
    db := setupTestDB(t)
    defer db.Close()

    repo := NewUserRepository(db)

    user := &models.User{
        Email:    "test@example.com",
        Phone:    "+6281234567890",
        FullName: "Test User",
    }

    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

### Integration Tests

```go
func TestUserFlow(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()

    // Initialize dependencies
    repo := repository.NewUserRepository(db)
    service := services.NewUserService(repo)

    // Test create
    user, err := service.CreateUser(ctx, &models.CreateUserRequest{
        Email:    "test@example.com",
        Phone:    "+6281234567890",
        FullName: "Test User",
        Password: "password123",
    })
    assert.NoError(t, err)

    // Test get
    fetchedUser, err := service.GetUser(ctx, user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Email, fetchedUser.Email)
}
```

## Migration

### Create New Migration

```bash
# Create up migration
cat > database/migrations/000002_add_users_profile.up.sql << 'EOF'
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);
ALTER TABLE users ADD COLUMN bio TEXT;
EOF

# Create down migration
cat > database/migrations/000002_add_users_profile.down.sql << 'EOF'
ALTER TABLE users DROP COLUMN IF EXISTS bio;
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
EOF
```

### Run Migration

```bash
make db-migrate-up
```

## Performance Optimization

### 1. Index Optimization

```sql
-- Add index for frequently queried columns
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone ON users(phone) WHERE deleted_at IS NULL;
```

### 2. Query Optimization

```go
// Use SELECT only needed columns
query := `SELECT id, email, full_name FROM users WHERE id = $1`

// Use LIMIT for pagination
query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
```

### 3. Connection Pooling

Adjust based on your load:
```go
db.SetMaxOpenConns(25)      // Max connections
db.SetMaxIdleConns(25)      // Idle connections
db.SetConnMaxLifetime(5 * time.Minute)
```

## Monitoring

### Query Logging

```go
// Enable query logging in development
if env == "development" {
    db = db.Unsafe() // Disable type checking for debugging
}
```

### Connection Stats

```go
stats := db.Stats()
logger.WithFields(logrus.Fields{
    "open_connections": stats.OpenConnections,
    "in_use":          stats.InUse,
    "idle":            stats.Idle,
    "wait_count":      stats.WaitCount,
    "wait_duration":   stats.WaitDuration,
}).Info("Database connection stats")
```

## Common Issues

### 1. Too Many Connections

**Solution**: Adjust MaxOpenConns or increase PostgreSQL max_connections

### 2. Connection Timeout

**Solution**: Increase context timeout or check network

### 3. Deadlock

**Solution**: Always acquire locks in the same order, use proper transaction isolation

### 4. Slow Queries

**Solution**: Add indexes, optimize queries, use EXPLAIN ANALYZE

## Resources

- [sqlx Documentation](https://jmoiron.github.io/sqlx/)
- [PostgreSQL Best Practices](https://wiki.postgresql.org/wiki/Don%27t_Do_This)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
