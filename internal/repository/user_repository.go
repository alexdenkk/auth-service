package repository

import (
	"alexdenkk/auth-service/internal/domain"
	"alexdenkk/auth-service/pkg/config"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

// User repository structure
type userRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// New user repository
func NewUserRepository(cfg *config.DBConfig, logger *zap.Logger) domain.UserRepository {
	// Connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	logger.Info("Connecting to database",
		zap.String("host", cfg.Host),
		zap.String("database", cfg.Name),
	)

	// Using postgres driver
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		logger.Error(
			"Failed to open database connection",
			zap.Error(err),
		)
		return nil
	}

	// Pinging database
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	logger.Info("Database connection established successfully")

	return &userRepository{
		logger: logger,
		db:     db,
	}
}

// Get user by email
func (repository *userRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	start := time.Now()

	query := `
		SELECT id, email, password, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`

	// Executing query
	var user domain.User

	err := repository.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	duration := time.Since(start)

	// Checking errors
	if err == sql.ErrNoRows {
		repository.logger.Warn("user not found by email",
			zap.String("layer", "repository"),
			zap.String("email", email),
		)
		return domain.User{}, err
	}

	if err != nil {
		repository.logger.Error("Failed to find user by email",
			zap.String("layer", "repository"),
			zap.String("email", email),
			zap.Error(err),
		)
		return domain.User{}, err
	}

	repository.logger.Info("user found by email successfully",
		zap.String("layer", "repository"),
		zap.String("email", email),
		zap.Duration("duration", duration),
	)

	return user, nil
}

func (repository *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	start := time.Now()

	query := `
		SELECT id, email, password, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`

	var user domain.User

	err := repository.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	duration := time.Since(start)

	if err == sql.ErrNoRows {
		repository.logger.Warn("user not found by id",
			zap.String("layer", "repository"),
			zap.String("id", id.String()),
		)
		return domain.User{}, err
	}

	if err != nil {
		repository.logger.Error("Failed to find user by id",
			zap.String("layer", "repository"),
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return domain.User{}, err
	}

	repository.logger.Info("user found by id successfully",
		zap.String("layer", "repository"),
		zap.String("id", id.String()),
		zap.Duration("duration", duration),
	)

	return user, nil
}

// User record creation
func (r *userRepository) CreateUser(ctx context.Context, user domain.User) error {
	start := time.Now()

	r.logger.Info("creating new user", zap.String("email", user.Email))

	query := `
		INSERT INTO users (id, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	// Executing query
	err := r.db.QueryRowContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	duration := time.Since(start)

	// Checking errors
	if err != nil {
		r.logger.Error("failed to create user",
			zap.String("layer", "repository"),
			zap.String("email", user.Email),
			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Info("user created successfully",
		zap.String("layer", "repository"),
		zap.String("email", user.Email),
		zap.String("user_id", user.ID.String()),
		zap.Duration("duration", duration),
	)

	return nil
}
