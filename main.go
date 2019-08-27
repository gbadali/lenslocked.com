package main

import (
	"github.com/gbadali/lenslocked.com/rand"

	"fmt"
	"net/http"

	"github.com/gorilla/csrf"

	"github.com/gbadali/lenslocked.com/controllers"
	"github.com/gbadali/lenslocked.com/middleware"
	"github.com/gbadali/lenslocked.com/models"
	"github.com/gorilla/mux"
)

// TODO: move to config file
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
	cfg := DefaultConfig()
	dbCfg := DefaultPostgresConfig()
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
	)

	if err != nil {
		panic(err)
	}

	defer services.Close()
	// Do a destructive reset on the DB for schema changes we can't migrate
	// services.DestructiveReset()
	services.AutoMigrate()

	// TODO: this could also be in the config file
	b, err := rand.Bytes(32)
	if err != nil {
		panic(err)
	}
	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	userMw := middleware.User{
		UserService: services.User,
	}

	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	newGallery := requireUserMw.Apply(galleriesC.New)
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)

	r.Handle("/", staticC.Home).
		Methods("GET")
	r.Handle("/contact", staticC.Contact).
		Methods("GET")
	r.Handle("/faq", staticC.FAQ).
		Methods("GET")
	// !User Routes
	r.HandleFunc("/signup", usersC.New).
		Methods("GET")
	r.HandleFunc("/signup", usersC.Create).
		Methods("POST")
	r.Handle("/login", usersC.LoginView).
		Methods("GET")
	r.HandleFunc("/login", usersC.Login).
		Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).
		Methods("GET")
	// !Gallery Routes
	r.Handle("/galleries",
		requireUserMw.ApplyFn(galleriesC.Index)).
		Methods("GET").
		Name(controllers.IndexGalleries)
	r.Handle("/galleries/new", newGallery).
		Methods("GET")
	r.HandleFunc("/galleries", createGallery).
		Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).
		Methods("GET").
		Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit",
		requireUserMw.ApplyFn(galleriesC.Edit)).
		Methods("GET").
		Name(controllers.EditGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/update",
		requireUserMw.ApplyFn(galleriesC.Update)).
		Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete",
		requireUserMw.ApplyFn(galleriesC.Delete)).
		Methods("POST")
	// !Image routes
	r.HandleFunc("/galleries/{id:[0-9]+}/images",
		requireUserMw.ApplyFn(galleriesC.ImageUpload)).
		Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete",
		requireUserMw.ApplyFn(galleriesC.ImageDelete)).
		Methods("POST")
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// !Other Routes
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", csrfMw(userMw.Apply(r)))
	// !Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
