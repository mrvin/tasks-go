package handlers //nolint:dupl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/pinger/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/pinger/pkg/http/response"
)

type HostCreator interface {
	CreateHost(ctx context.Context, host *storage.Host) error
}

func NewCreateHost(creator HostCreator) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Read json request
		var request storage.Host

		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := creator.CreateHost(req.Context(), &request); err != nil {
			err := fmt.Errorf("create host: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res, http.StatusOK)

		slog.Info("Host creation was successful")
	}
}
