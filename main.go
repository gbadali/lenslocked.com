package main

import (
	"fmt"
	"net/http"
<<<<<<< HEAD

	"github.com/julienschmidt/httprouter"
)

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Welcome to my site!</h1>")
}

func contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "to get in touch, please send an email "+
		"to <a href=\"mailto:support@lenslocked.com\">"+
		"support@lenslockedd.com</a>.")
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
	router := httprouter.New()
	router.GET("/", home)
	router.GET("/contact", contact)
	router.GET("/faq", faq)
	var h http.Handler = http.HandlerFunc(notFoundHandler)
	router.NotFound = h
	http.ListenAndServe(":3000", router)
=======
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Welcome to my site</h1>")
	} else if r.URL.Path == "/contact" {
		fmt.Fprint(w, "to get in touch, please send an email "+
			"to <a href=\"mailto:support@lenslocked.com\">"+
			"support@lenslockedd.com</a>.")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>We could not find the page "+
			"were looking for :(</h1>"+
			"<p>Please email us if you keep being sent to an "+
			"invalid page.</p>")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
>>>>>>> parent of b7446cd... 3.4.2 Ex2 - add custom 404 page
}
