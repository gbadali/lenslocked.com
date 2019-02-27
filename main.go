package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var homeTemplate *template.Template
var contactTemplate *template.Template

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func faq(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Freqently Asked Questions</h1>")
	fmt.Fprint(w, "<ul>"+
		"<li> Why is HTML so annoying</li>"+
		"<li> This is an un-ordered list</li>"+
		"<li> It should have bullet points</li>"+
		"<li> I should probably move my html to a separate folder</li></ul>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound) // StatusNotFound = 404
	fmt.Fprint(w, "Couldn't find page 404&#128169")
}

func main() {
	var err error
	homeTemplate, err = template.ParseFiles("views/home.html")
	if err != nil {
		panic(err)
	}
	contactTemplate, err = template.ParseFiles("views/contact.html")
	if err != nil {
		panic(err)
	}
	router := httprouter.New()
	router.GET("/", home)
	router.GET("/contact", contact)
	router.GET("/faq", faq)
	var h http.Handler = http.HandlerFunc(notFoundHandler)
	router.NotFound = h
	http.ListenAndServe(":3000", router)
}
