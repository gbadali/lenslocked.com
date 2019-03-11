package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "development"
	dbname   = "lenselocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("succesfully connected!")
	defer db.Close()

	_, err = db.Exec(`
	INSERT INTO users(name, email)
	VALUES ($1, $2)`,
		"Jon Calhoun", "jon@calhoun.io")
	if err != nil {
		panic(err)
	}

	var id int
	row := db.QueryRow(`
	INSERT INTO users(name, email)
	VALUES ($1, $2) RETURNING id`,
		"Jon Calhoun", "jon@calhooun.io")
	err = row.Scan(&id)
	if err != nil {
		panic(err)
	}
}
