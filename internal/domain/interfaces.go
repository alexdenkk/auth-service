package domain

import (
	"context"
	"github.com/google/uuid"
)

// User repository interface
type UserRepository interface {
	GetUserByEmail(context.Context, string) (User, error)
	GetUserByID(context.Context, uuid.UUID) (User, error)
	CreateUser(context.Context, User) error
}

// Service interface
type Service interface {
	SignIn(context.Context, string, string) (string, error)
	SignUp(context.Context, string, string) error
	GetSelf(context.Context, string) (User, error)
}
