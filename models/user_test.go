package models_test

import (
	"testing"

	"github.com/ahugues/porygon-backend-go/models"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Parallel()
	user, err := models.NewUser("user", "password", "John", "Doe", "test@example.com")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	assert.Equal(t, user.FirstName, "John", "First Name")
	assert.Equal(t, user.LastName, "Doe", "Last Name")
	assert.Equal(t, user.Login, "user", "Login")
	assert.Equal(t, user.Email, "test@example.com", "Email")
	assert.Equal(t, user.CheckPasswd("password"), true)
	assert.Equal(t, user.CheckPasswd("password2"), false)
}
