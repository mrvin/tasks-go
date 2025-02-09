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

type PingCreator interface {
	CreatePing(ctx context.Context, ping *storage.Ping) error
}

func NewCreatePing(creator PingCreator) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Read json request
		var request storage.Ping

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

		if err := creator.CreatePing(req.Context(), &request); err != nil {
			err := fmt.Errorf("create ping: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res, http.StatusOK)

		slog.Info("Ping creation was successful")
	}
}
