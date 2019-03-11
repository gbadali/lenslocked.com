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

	var id int
	for i := 1; i < 6; i++ {
		// create somoe fake data
		userID := 1
		if i > 3 {
			userID = 2
		}
		amount := 1000 * i
		description := fmt.Sprintf("USB-C Adapter x%d", i)

		err = db.QueryRow(`
			INSERT INTO orders (user_id, amount, description)
			VALUES ($1, $2, $3)
			RETURNING id`,
			userID, amount, description).Scan(&id)
		if err != nil {
			panic(err)
		}
		fmt.Println("Created an ordder with the ID:", id)
	}
}
