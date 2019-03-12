package main

import (
	"fmt"

	"github.com/gbadali/lenslocked.com/models"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "development"
	dbname   = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	foundUser, err := us.ByEmail("michael@dundermifflin.com")
	if err != nil {
		panic(err)
	}

	if err := us.Delete(foundUser.ID); err != nil {
		panic(err)
	}
	// Verify the user is deleted
	_, err = us.ByID(foundUser.ID)
	if err != models.ErrNotFound {
		panic("user was not deleted")
	}
}
