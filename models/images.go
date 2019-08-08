package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

// Image is used to represent images stored in a Gallery.
// Image is NOT stored in the database, and instead
// refereences data storedd on disk.
type Image struct {
	GalleryID uint
	Filename  string
}

// Path is used to build the absolute path used to reference this image
// via a web request.
func (i *Image) Path() string {
	return "/" + i.RelativePath()
}

// RelativePath is used to build the path to this image on our local
// disk, relative to where our Go application is run from.
func (i *Image) RelativePath() string {
	// Convert the gallery ID to a string
	galleryID := fmt.Sprintf("%v", i.GalleryID)
	return filepath.ToSlash(filepath.Join("images", "galleries", galleryID,
		i.Filename))
}

func (is *imageService) Create(galleryID uint,
	r io.Reader, filename string) error {
	path, err := is.mkIMagePath(galleryID)
	if err != nil {
		return err
	}
	// Create a destination file
	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy the reader data to the destination file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}
	// Setup the Image slice we are returning
	ret := make([]Image, len(strings))
	for i, imgStr := range strings {
		ret[i] = Image{
			Filename:  filepath.Base(imgStr),
			GalleryID: galleryID,
		}
	}
	return ret, nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return filepath.Join("images", "galleries",
		fmt.Sprintf("%v", galleryID))
}

func (is *imageService) mkIMagePath(galleryID uint) (string, error) {
	// Setup to create files to store the gallery
	galleryPath := is.imagePath(galleryID)
	// Create our directory (an any necessary parent dirs)
	// using 0755 permissions.
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil

}
