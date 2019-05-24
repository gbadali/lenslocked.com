package models

import (
	"github.com/jinzhu/gorm"
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
type GalleryDB interface {
	Create(gallery *Gallery) error
}

type galleryService struct {
	GalleryDB
}
type GalleryValidator struct {
	GalleryDB
}

// check to make sure galleryGorm impements the GalleryDB methods
var _ GalleryDB = &galleryGorm{}

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
