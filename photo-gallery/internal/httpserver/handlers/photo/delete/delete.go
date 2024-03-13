package delete

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	httpresponse "github.com/mrvin/tasks-go/photo-gallery/pkg/http/response"
)

type PhotoDeleter interface {
	DeletePhoto(ctx context.Context, name string) error
}

func New(photoDeleter PhotoDeleter, dirPhotos string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		namePhoto := req.URL.Query().Get("name")
		if namePhoto == "" {
			err := errors.New("Empty file name")
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := os.Remove(filepath.Join(dirPhotos, namePhoto)); err != nil {
			err := fmt.Errorf("Delete photo: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		ext := filepath.Ext(namePhoto) // e.g., ".jpg", ".JPEG"
		nameThumbnail := strings.TrimSuffix(namePhoto, ext) + ".thumb" + ext
		if err := os.Remove(filepath.Join(dirPhotos, nameThumbnail)); err != nil {
			err := fmt.Errorf("Delete photo: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := photoDeleter.DeletePhoto(req.Context(), namePhoto); err != nil {
			err := fmt.Errorf("Delete photo info: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res)

		slog.Info("Photo removed", slog.String("name", namePhoto))
	}
}
