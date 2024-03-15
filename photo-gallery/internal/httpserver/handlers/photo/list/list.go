package list

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/photo-gallery/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/photo-gallery/pkg/http/response"
)

type PhotoLister interface {
	ListPhotos(ctx context.Context) ([]storage.PhotoInfo, error)
}

type ResponseListPhotos struct {
	ListPhotos []storage.PhotoInfo `json:"listPhotos"`
	Status     string              `json:"status"`
}

func New(photoLister PhotoLister) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		list, err := photoLister.ListPhotos(req.Context())
		if err != nil {
			err := fmt.Errorf("get list photo info: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		}
		// Write json response
		response := ResponseListPhotos{
			ListPhotos: list,
			Status:     "OK",
		}

		jsonResponseListPhotos, err := json.Marshal(response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponseListPhotos); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
