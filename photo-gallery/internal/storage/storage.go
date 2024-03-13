package storage

import (
	"context"
)

type PhotoInfo struct {
	Name         string `json:"name"`
	URLPhoto     string `json:"urlPhoto"`
	URLThumbnail string `json:"urlThumbnail"`
}

type PhotoStorage interface {
	SavePhoto(ctx context.Context, photoInfo *PhotoInfo) error
	ListPhotos(ctx context.Context) ([]PhotoInfo, error)
	DeletePhoto(ctx context.Context, name string) error
}
