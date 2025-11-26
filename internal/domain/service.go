package domain

import (
	"alexdenkk/auth-service/pkg/config"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Service structure
type service struct {
	userRepository UserRepository
	jwtConfig      *config.JwtConfig
	logger         *zap.Logger
}

// New service
func NewService(cfg *config.JwtConfig, log *zap.Logger, repo UserRepository) Service {
	return &service{
		jwtConfig:      cfg,
		logger:         log,
		userRepository: repo,
	}
}

// Register user
func (service *service) SignUp(ctx context.Context, email, password string) error {
	start := time.Now()

	service.logger.Info("processing registration request",
		zap.String("layer", "domain"),
		zap.String("email", email),
	)

	// Checking if email already exists
	_, err := service.userRepository.GetUserByEmail(ctx, email)

	if err == nil {
		service.logger.Warn("user with email already registered",
			zap.String("layer", "domain"),
			zap.Error(err),
		)

		return errors.New("user with this email already exists")
	}

	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		service.logger.Warn("error hashing user password",
			zap.String("layer", "domain"),
			zap.Error(err),
		)

		return errors.New("error while registering user")
	}

	// User struct
	user := User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
	}

	// User validation
	if err := user.Validate(); err != nil {
		service.logger.Warn("error validating user",
			zap.String("layer", "domain"),
			zap.Error(err),
		)

		return err
	}

	// Creating user record
	err = service.userRepository.CreateUser(ctx, user)

	if err != nil {
		service.logger.Warn("error creating user",
			zap.String("layer", "domain"),
			zap.Error(err),
		)

		return errors.New("error while registering user")
	}

	service.logger.Info("successfully registered user",
		zap.String("layer", "domain"),
		zap.String("email", user.Email),
		zap.Duration("duration", time.Since(start)),
	)

	return nil
}

// User authorization
func (service *service) SignIn(ctx context.Context, email, password string) (string, error) {
	start := time.Now()

	service.logger.Info("processing authorization request",
		zap.String("layer", "domain"),
		zap.String("email", email),
	)

	// Getting user record by email
	user, err := service.userRepository.GetUserByEmail(ctx, email)

	if err != nil {
		service.logger.Warn("user not found",
			zap.String("layer", "domain"),
			zap.Error(err),
		)

		return "", errors.New("user not found")
	}

	// Checking password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		service.logger.Warn("invalid user password",
			zap.String("layer", "domain"),
			zap.Error(err),
		)

		return "", errors.New("invalid password")
	}

	// Token claims
	claims := &JwtClaims{
		ID:    user.ID,
		Email: user.Email,
	}

	// Generating token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tkn, err := token.SignedString(service.jwtConfig.SignKey)

	if err != nil {
		service.logger.Warn("error generating token",
			zap.String("layer", "domain"),
			zap.Error(err),
		)

		return "", errors.New("authorization error")
	}

	service.logger.Info("successfully authorized user",
		zap.String("layer", "domain"),
		zap.String("email", user.Email),
		zap.Duration("duration", time.Since(start)),
	)

	return "Bearer " + tkn, nil
}

// Get user by his token
func (service *service) GetSelf(ctx context.Context, token string) (User, error) {
	start := time.Now()

	service.logger.Info("processing get user self request",
		zap.String("layer", "domain"),
	)

	// Splitting token
	parts := strings.Split(token, " ")

	// Checking token
	if len(parts) != 2 || parts[0] != "Bearer" {
		service.logger.Warn("invalid token",
			zap.String("layer", "domain"),
		)

		return User{}, errors.New("invalid authorization token")
	}

	// Parsing claims
	var claims JwtClaims

	parsed, err := jwt.ParseWithClaims(
		parts[1],
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return service.jwtConfig.SignKey, nil
		},
	)

	// Checking errors
	if err != nil {
		service.logger.Warn("token parsing error",
			zap.String("layer", "domain"),
		)
		return User{}, errors.New("invalid authorization token")
	}

	// Checking is parsing valid
	if !parsed.Valid {
		service.logger.Warn("token parsed invalid",
			zap.String("layer", "domain"),
		)
		return User{}, errors.New("invalid authorization token")
	}

	// Getting user record
	user, err := service.userRepository.GetUserByID(ctx, claims.ID)

	if err != nil {
		service.logger.Warn("error getting user",
			zap.String("layer", "domain"),
			zap.Error(err),
		)

		return User{}, errors.New("error getting user")
	}

	service.logger.Info("successfully authorized user",
		zap.String("layer", "domain"),
		zap.String("email", user.Email),
		zap.Duration("duration", time.Since(start)),
	)

	return user, nil
}
