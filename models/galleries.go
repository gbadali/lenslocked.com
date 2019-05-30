package models

import (
	"github.com/jinzhu/gorm"
)

const (
	ErrUserIDRequired modelError = "models: user ID is required"
	ErrTitleRequired  modelError = "models: title is required"
)

type Gallery struct {
	gorm.Model
	// a unique id for each user
	UserID uint `gorm:"not_null;index`
	// the title of each gallery
	Title string `gorm:"not_null`
}
type GalleryService interface {
	GalleryDB
}

//GalleryDB is used to interact with the galleries database.
//
// For pretty much all single gallery queries:
// If the gallery is found, we will return a nil error
// If the gallery is not found, we will return ErrNotFound
// If there is another error, we will return an eror with
// more information about what went wrong.  This may not be
// an error generated by the models package
type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
}

type galleryService struct {
	GalleryDB
}
type GalleryValidator struct {
	GalleryDB
}

func (gv *GalleryValidator) userIDRequired(g *Gallery) error {
	if g.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *GalleryValidator) TitleRequired(g *Gallery) error {
	if g.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (gv *GalleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFns(gallery,
		gv.userIDRequired,
		gv.TitleRequired)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

// check to make sure galleryGorm impements the GalleryDB methods
var _ GalleryDB = &galleryGorm{}

type galleryValFn func(*Gallery) error

func runGalleryValFns(gallery *Gallery, fns ...galleryValFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

type galleryGorm struct {
	db *gorm.DB
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &GalleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}
