package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWT token structure
type JwtClaims struct {
	jwt.RegisteredClaims

	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
