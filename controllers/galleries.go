package controllers

import (
	"log"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/gbadali/lenslocked.com/context"
	"github.com/gbadali/lenslocked.com/models"
	"github.com/gbadali/lenslocked.com/views"
)

// Named routes
const (
	IndexGalleries = "index_galleries"
	ShowGallery    = "show_gallery"
	EditGallery    = "edit_gallery"
)

// Other consts
const (
	maxMultipartMem = 1 << 20 // 1 megabyte
)

// GalleryForm provides some structure to the new gallery form
type GalleryForm struct {
	Title string `schema: "title"`
}

// NewGalleries instantiates the Gallery object and adds the templates
func NewGalleries(gs models.GalleryService, is models.ImageService, r *mux.Router) *Galleries {
	return &Galleries{
		New:       views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		gs:        gs,
		is:        is,
		r:         r,
	}
}

// Galleries stores the different views that the Gallery has
type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	is        models.ImageService
	r         *mux.Router
}

// Create takes the information from the new Gallery form and uses
// it to set up a Gallery
// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}
	user := context.User(r.Context())
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}

	url, err := g.r.Get(EditGallery).URL("id",
		strconv.Itoa(int(gallery.ID)))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// Show ...
// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, r, vd)
}

// Edit ...
// GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	// A user needs to be logged in to access this page, so we can
	// assume that the RequreUser middleware ahs run and
	// set tht user for us in the request context.
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, r, vd)
}

// Update ...
// POST /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in controllers/galleries.go:Update function")
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		// If there is an error we are going to render the
		// EditView again with an alert message
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	// If there is an error our alert will be an error. Otherwise
	// we will still rendder an alert, but insteadd it will be
	// a success message.
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Gallery successfully updated!",
		}
	}
	g.EditView.Render(w, r, vd)
}

// Delete ...
// POST /galleries/:idd/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	// Loookup the gallery using the galleryByID we wrote earlier
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		// If there is an error the galleryByID will have renderd
		// it for us already
		return
	}

	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		// If there is an error we wnat to set an alert and
		// render the edit page with the error.  We also need
		// to set the Yield to gallery so that the EditView
		// is rendered correctly.
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, r, vd)
		return
	}
	url, err := g.r.Get(IndexGalleries).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// POST /galleries/:id/images/:filename/delete
func (g *Galleries) ImageDelete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permisiont to edit "+
			"this gallery or image", http.StatusForbidden)
		return
	}
	// Get the filename from the path
	filename := mux.Vars(r)["filename"]
	// Build the Image model
	i := models.Image{
		Filename:  filename,
		GalleryID: gallery.ID,
	}
	// Try to delete the image.
	err = g.is.Delete(&i)
	if err != nil {
		// Render the edit page with any errors.
		var vd views.Data
		vd.Yield = gallery
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	// If all goes well, redirect to the edit gallery page.
	url, err := g.r.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// Index ...
// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
}

// ImageUpload ...
// POST /galleries/:id/images
func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	// Parse the multipart form, the results go back into r
	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	// Iterate over uploaded files to process them.
	files := r.MultipartForm.File["images"]
	for _, f := range files {
		// Open the uploaded file
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer file.Close()
		err = g.is.Create(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
	}

	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Images successfully uploaded!",
	}
	g.EditView.Render(w, r, vd)
}

func (g *Galleries) galleryByID(w http.ResponseWriter,
	r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)

		}
		return nil, err
	}
	images, _ := g.is.ByGalleryID(gallery.ID)
	gallery.Images = images
	return gallery, nil
}


