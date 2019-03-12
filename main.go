package main

import (
	"fmt"
	"net/http"

	"github.com/gbadali/lenslocked.com/models"

	"github.com/gorilla/mux"

	"github.com/gbadali/lenslocked.com/controllers"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "develpment"
	dbname   = "lenslocked_dev"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound) // StatusNotFound = 404
	fmt.Fprint(w, "Couldn't find page 404&#128169")
}

func main() {
	// Create a DB connection string and then use it to
	// create our model services.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.FAQ).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
