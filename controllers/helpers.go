package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	dec := schema.NewDecoder()
	// Call the IgnoreUnknownKeys function to tell schema's decoder
	// to ignore the CSRF token key
	dec.IgnoreUnknownKeys(true)
	if err := dec.Decode(dst, r.PostForm); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
