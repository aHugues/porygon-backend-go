package services_test

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"os"
	"testing"

	porygon_errors "github.com/ahugues/porygon-backend-go/errors"
	"github.com/ahugues/porygon-backend-go/models"
	"github.com/ahugues/porygon-backend-go/services"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
)

var porygonTestDB string
var dbDriver string

func TestMain(m *testing.M) {
	var dbType string
	flag.StringVar(&dbType, "dbtype", "postgres", "Database driver (postgres or mysql)")
	flag.Parse()

	log.Printf("coucou")

	if dbType == "postgres" {
		porygonTestDB = "postgres://azure:asure@localhost/porygontest"
		dbDriver = "pgx"
	} else if dbType == "mysql" {
		porygonTestDB = "mysql://toto"
		dbDriver = "mysql"
	} else {
		log.Fatalf("Unexpected dbtype %s, expected mysql or postgres", dbType)
	}

	os.Exit(m.Run())
}

func TestGetUserOK(t *testing.T) {
	t.Parallel()
	db, err := sql.Open(dbDriver, porygonTestDB)
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}
	defer db.Close()

	s := services.NewConcreteUserService(db)
	if usr, err := s.GetUser("user1"); err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	} else {
		assert.Equal(t, "user1", usr.Login)
		assert.Equal(t, "9f2b6641-fbdb-4671-b52b-087141257116", usr.UUID.String())
		assert.Equal(t, "User", usr.FirstName)
		assert.Equal(t, "One", usr.LastName)
		assert.Equal(t, "$2a$10$mhz0tYmWYvHZkFPQgu5nVeyVU07h82RDS9GmDtVZOSyPIFZ7IHHPK", usr.Password)
		assert.Equal(t, "user1@email.com", usr.Email)
	}
}

func TestGetUserNOK(t *testing.T) {
	t.Parallel()
	db, err := sql.Open(dbDriver, porygonTestDB)
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}
	defer db.Close()

	s := services.NewConcreteUserService(db)
	if _, err := s.GetUser("no_user"); err == nil {
		t.Error("Unexpected nil error")
	} else if err.Error() != "Error getting user no_user: sql: no rows in result set" {
		t.Errorf("Unexpected error. Expected %s | got %s", "Error getting user no_user: sql: no rows in result set", err.Error())
	}
}

func TestGetUserInvalidUUID(t *testing.T) {
	t.Parallel()
	db, err := sql.Open(dbDriver, porygonTestDB)
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}
	defer db.Close()

	s := services.NewConcreteUserService(db)
	if _, err := s.GetUser("invalid_uuid"); err == nil {
		t.Error("Unexpected nil error")
	} else if err.Error() != "Invalid uuid 9f2b6641-fbdb-4671-b52b: invalid UUID length: 23" {
		t.Errorf("Unexpected error. Expected %s | got %s", "Invalid uuid 9f2b6641-fbdb-4671-b52b: invalid UUID length: 23", err.Error())
	}
}

func TestSaveUserOK(t *testing.T) {
	t.Parallel()
	db, err := sql.Open(dbDriver, porygonTestDB)
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}
	defer db.Close()

	s := services.NewConcreteUserService(db)
	newUser, err := models.NewUser("new_user_1", "toto", "New user", "One", "new1@email.com")
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}
	if err := s.SaveUser(newUser); err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	testedUser, err := s.GetUser("new_user_1")
	if err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	} else {
		assert.Equal(t, newUser, testedUser)
	}
}

func TestLoginOK(t *testing.T) {
	t.Parallel()
	db, err := sql.Open(dbDriver, porygonTestDB)
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}
	defer db.Close()

	s := services.NewConcreteUserService(db)
	if info, err := s.CheckLogin("user1", "password1"); err != nil {
		t.Errorf("Unexpected error %s", err.Error())
	} else {
		assert.Equal(t, "User", info.FirstName)
		assert.Equal(t, "One", info.LastName)
		assert.Equal(t, "user1", info.Login)
		assert.NotEqual(t, "", info.Token)
	}
}

func TestLoginNOK(t *testing.T) {
	t.Parallel()
	db, err := sql.Open(dbDriver, porygonTestDB)
	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}
	defer db.Close()

	s := services.NewConcreteUserService(db)
	if _, err := s.CheckLogin("user1", "password2"); !errors.Is(err, porygon_errors.ErrInvalidLogin) {
		t.Errorf("Unexpected error %s", err.Error())
	}
}
