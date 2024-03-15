//nolint:stylecheck
package save

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mrvin/tasks-go/photo-gallery/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/photo-gallery/pkg/http/response"
	"github.com/mrvin/tasks-go/photo-gallery/pkg/thumbnail"
)

type PhotoSaver interface {
	SavePhoto(ctx context.Context, photoInfo *storage.PhotoInfo) error
}

func New(photoSaver PhotoSaver, dirPhotos, addr, path string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		namePhoto := req.URL.Query().Get("name")
		if namePhoto == "" {
			err := errors.New("Empty file name")
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		photoPath := filepath.Join(dirPhotos, namePhoto)
		filePhoto, err := os.Create(photoPath)
		if err != nil {
			err := fmt.Errorf("Create photo file: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		defer filePhoto.Close()

		photoSize, err := io.Copy(filePhoto, req.Body)
		if err != nil {
			err := fmt.Errorf("Copy request body to file: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		defer req.Body.Close()

		nameThumbnail, err := thumbnail.ImageFile(photoPath)
		if err != nil {
			err := fmt.Errorf("Create thumbnail: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		_, nameThumbnail = filepath.Split(nameThumbnail)
		urlPhoto := "http://" + addr + path + "/" + namePhoto
		urlThumbnail := "http://" + addr + path + "/" + nameThumbnail
		photoInfo := &storage.PhotoInfo{
			Name:         namePhoto,
			URLPhoto:     urlPhoto,
			URLThumbnail: urlThumbnail}
		if err := photoSaver.SavePhoto(req.Context(), photoInfo); err != nil {
			err := fmt.Errorf("Save photo info: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res, http.StatusCreated)

		slog.Info("Photo upload", slog.String("name", namePhoto), slog.Int64("bytes", photoSize))
	}
}
