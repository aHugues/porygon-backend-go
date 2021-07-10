package services

import (
	"database/sql"
	"fmt"

	porygon_errors "github.com/ahugues/porygon-backend-go/errors"
	"github.com/ahugues/porygon-backend-go/models"
	"github.com/ahugues/porygon-backend-go/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type jsonUser struct {
	UUID      string `json:"uuid"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

type UserService interface {
	GetUser(id string) (*models.User, error)
	SaveUser(user *models.User) error
	CheckLogin(login string, password string) (*models.AuthInfo, error)
}

type ConcreteUserService struct {
	db *sql.DB
}

func NewConcreteUserService(db *sql.DB) UserService {
	return &ConcreteUserService{
		db: db,
	}
}

// GetUser gets the information for a user with the given login
func (s *ConcreteUserService) GetUser(id string) (user *models.User, err error) {
	jsonUser := new(jsonUser)
	if err := s.db.QueryRow(`SELECT uuid, login, password, "firstName", "lastName", email FROM "User" WHERE login = $1`, id).Scan(&jsonUser.UUID, &jsonUser.Login, &jsonUser.Password, &jsonUser.FirstName, &jsonUser.LastName, &jsonUser.Email); err != nil {
		return nil, fmt.Errorf("Error getting user %s: %w", id, err)
	}
	user = &models.User{
		Login:     jsonUser.Login,
		Password:  jsonUser.Password,
		FirstName: jsonUser.FirstName,
		LastName:  jsonUser.LastName,
		Email:     jsonUser.Email,
	}
	parsedUUID, err := uuid.Parse(jsonUser.UUID)
	if err != nil {
		return nil, fmt.Errorf("Invalid uuid %s: %w", jsonUser.UUID, err)
	}
	user.UUID = parsedUUID
	return
}

// SaveUser creates the new user
func (s *ConcreteUserService) SaveUser(user *models.User) error {
	tsx, err := s.db.Begin()
	if err != nil {
		tsx.Rollback()
		return fmt.Errorf("Error starting transaction: %w", err)
	}

	stmt, err := tsx.Prepare(`INSERT INTO "User"(uuid, login, password, "firstName", "lastName", email) VALUES($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		tsx.Rollback()
		return fmt.Errorf("Error preparing query: %w", err)
	}

	_, err = stmt.Exec(user.UUID.String(), user.Login, user.Password, user.FirstName, user.LastName, user.Email)
	if err != nil {
		tsx.Rollback()
		return fmt.Errorf("Error executing query: %w", err)
	}
	return tsx.Commit()
}

// CheckLogin verifies that the provided login and password correspond to an actual user. If this is the case, a new token
// will be provided for further authentication
func (s *ConcreteUserService) CheckLogin(login string, password string) (user *models.AuthInfo, err error) {
	fullUser, err := s.GetUser(login)
	if err != nil {
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(fullUser.Password), []byte(password)); err != nil {
		return nil, porygon_errors.ErrInvalidLogin
	}
	wrapper := utils.JwtWrapper{
		SecretKey:       "toto",
		Issuer:          "porygon",
		ExpirationHours: 24,
	}
	token, err := wrapper.GenerateToken(fullUser.UUID.String())
	if err != nil {
		return
	}
	user = models.NewAuthInfo(fullUser.Login, fullUser.FirstName, fullUser.LastName, fullUser.Email, token)
	return
}
