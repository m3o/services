package photos

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strings"

	cloudinary "github.com/ben-toogood/go-cloudinary"
	uuid "github.com/satori/go.uuid"
)

// ErrCredentials is returned when the creadentials provided are invalid
var ErrCredentials = errors.New("Invalid cloudinary credentials")

// ErrInvalidData is returned when an invalid Base64 image is provided
var ErrInvalidData = errors.New("Invalid image provided")

// ErrUnknown is returned when an unknown error occurred
var ErrUnknown = errors.New("Unknown error occurred")

// New returns an instance of the photos service
func New(creds string) (Service, error) {
	cloud, err := cloudinary.Dial(creds)

	if err != nil {
		return Service{}, ErrCredentials
	}

	return Service{cloud}, nil
}

// Service is a instance of an image upload service
type Service struct {
	cloud *cloudinary.Service
}

// Upload takes a base64 encoded photo, uploads it and returns the UUID
func (s Service) Upload(input string) (string, error) {
	base64Comps := strings.Split(input, ",")
	if len(base64Comps) > 2 {
		return "", ErrInvalidData
	}
	str := base64Comps[len(base64Comps)-1]

	base64Bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", ErrInvalidData
	}

	name := uuid.NewV4().String()

	if _, err := s.cloud.UploadImage(name, bytes.NewReader(base64Bytes), ""); err != nil {
		return "", err
	}

	return name, nil
}

// GetURL returns the URL for the photo. Size can be provided, width then height in pixels.
func (s Service) GetURL(uuid string, size ...int) string {
	if uuid == "" {
		return ""
	}

	var width, height = 200, 200

	switch len(size) {
	case 1:
		width = size[0]
	case 2:
		width = size[0]
		height = size[1]
	}

	return s.cloud.TransformedImageURL(uuid, cloudinary.SizeTransformation{
		Height: height, Width: width,
	})
}
