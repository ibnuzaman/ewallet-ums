package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/ibnuzaman/ewallet-ums/helpers"
	"github.com/ibnuzaman/ewallet-ums/internal/constants"
)

var (
	oncePostgres sync.Once
	pgSqlx       *sqlx.DB
)

// PostgresConfig holds the database configuration.
type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// GetPostgresConfig returns database configuration from environment variables.
func GetPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:            helpers.GetEnv("DB_HOST", "localhost"),
		Port:            helpers.GetEnv("DB_PORT", "5432"),
		User:            helpers.GetEnv("DB_USER", "postgres"),
		Password:        helpers.GetEnv("DB_PASSWORD", "postgres"),
		DBName:          helpers.GetEnv("DB_NAME", "ewallet_ums"),
		SSLMode:         helpers.GetEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    constants.DefaultMaxOpenConns,
		MaxIdleConns:    constants.DefaultMaxIdleConns,
		ConnMaxLifetime: constants.DefaultConnMaxLifetime,
		ConnMaxIdleTime: constants.DefaultConnMaxIdleTime,
	}
}

// InitPostgres initializes PostgreSQL connection.
func InitPostgres() (*sqlx.DB, error) {
	var err error

	oncePostgres.Do(func() {
		config := GetPostgresConfig()

		// Build DSN
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
			config.SSLMode,
		)

		autoMigrate(dsn)

		// Open connection
		pgSqlx, err = sqlx.Open("postgres", dsn)
		if err != nil {
			helpers.Logger.Errorf("Failed to open database connection: %v", err)
			return
		}

		// Set connection pool settings
		pgSqlx.SetMaxOpenConns(config.MaxOpenConns)
		pgSqlx.SetMaxIdleConns(config.MaxIdleConns)
		pgSqlx.SetConnMaxLifetime(config.ConnMaxLifetime)
		pgSqlx.SetConnMaxIdleTime(config.ConnMaxIdleTime)

		// Test connection with timeout
		ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultPingTimeout)
		defer cancel()

		err = pgSqlx.PingContext(ctx)
		if err != nil {
			helpers.Logger.Errorf("Failed to ping database: %v", err)
			return
		}

		helpers.Logger.Info("Successfully connected to PostgreSQL database")
	})

	return pgSqlx, err
}

// GetPostgresDB returns the PostgreSQL database instance.
func GetPostgresDB() *sqlx.DB {
	return pgSqlx
}

// ClosePostgres closes the PostgreSQL connection.
func ClosePostgres() error {
	if pgSqlx != nil {
		helpers.Logger.Info("Closing PostgreSQL connection...")
		return pgSqlx.Close()
	}
	return nil
}

// HealthCheck checks if the database connection is healthy.
func HealthCheck(ctx context.Context) error {
	if pgSqlx == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	if err := pgSqlx.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

func autoMigrate(dsn string) {
	baseDir := "database/migrations"
	files, err := os.ReadDir(baseDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			helpers.Logger.Warnf("Migration directory does not exist: %s", baseDir)
			return
		}
		helpers.Logger.Errorf("Failed to read migration directory: %v", err)
		return
	}

	if len(files) == 0 {
		helpers.Logger.Warnf("No migration files found in %s", baseDir)
		return
	}

	m, err := migrate.New("file://"+baseDir, dsn)
	if err != nil {
		helpers.Logger.Errorf("Failed to create migration instance: %v", err)
		return
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		helpers.Logger.Errorf("Migration failed: %v", err)
		return
	}

	helpers.Logger.Infof("Found %d migration files in %s", len(files), baseDir)
}
