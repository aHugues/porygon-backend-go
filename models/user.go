package models

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user for authentication purposes
type User struct {
	UUID      uuid.UUID `json:"uuid"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
}

// NewUser creates a new User with the given information
func NewUser(login string, password string, firstName string, lastName string, email string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error generating password hash: %w", err)
	}
	return &User{
		UUID:      uuid.New(),
		Password:  string(hash),
		Login:     login,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}, nil
}

// CheckPasswd returns true if the provided password matches the hashed password
func (u *User) CheckPasswd(passwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwd)) == nil
}

// AuthInfo represents the authentication information for a user to establish a session
type AuthInfo struct {
	Login     string `json:"login"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Token     string `json:"token"`
}

// NewAuthInfo returns a new AuthInfo struct from the given data
func NewAuthInfo(login string, firstName string, lastName string, email string, token string) *AuthInfo {
	return &AuthInfo{
		Login:     login,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Token:     token,
	}
}
