package models

import (
	"github.com/jinzhu/gorm"
)

type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index`
	Title  string `gorm:"not_null`
}

type GalleryService interface {
	Create(gallery *Gallery) error
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	// TODO: Implement this later
	return nil
}
