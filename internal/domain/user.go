package domain

import (
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// User structure
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hash
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User structure validation
func (u *User) Validate() error {
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid user email")
	}

	return nil
}
