package controllers

import (
	"fmt"
	"net/http"

	"github.com/gbadali/lenslocked.com/context"
	"github.com/gbadali/lenslocked.com/models"
	"github.com/gbadali/lenslocked.com/views"
)

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

type GalleryForm struct {
	Title string `schema: "title"`
}

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	user := context.User(r.Context())
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
