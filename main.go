package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ahugues/porygon-backend-go/config"
	"github.com/ahugues/porygon-backend-go/models"
	"github.com/ahugues/porygon-backend-go/services"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	testConf := &config.DatabaseConfig{
		Driver:   "pgx",
		User:     "azure",
		Password: "porygon",
		Host:     "localhost",
		Port:     5432,
		Database: "porygontest",
	}

	db, err := sql.Open(testConf.Driver, testConf.ToEndpoint())
	if err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	defer db.Close()

	userService := services.NewConcreteUserService(db)
	user, err := models.NewUser("user2", "password2", "test", "user", "test@email.com")
	if err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	if err := userService.SaveUser(user); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}

	authInfo, err := userService.CheckLogin("porygon", "porygon")
	if err != nil {
		log.Fatalf("Unexpecter error %s", err.Error())
	}
	fmt.Println(authInfo)
}
