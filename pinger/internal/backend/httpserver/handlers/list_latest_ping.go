package handlers //nolint:dupl

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/pinger/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/pinger/pkg/http/response"
)

type LatestPingLister interface {
	ListLatestPing(ctx context.Context) ([]storage.Ping, error)
}

//nolint:tagliatelle
type ResponseListLatestPing struct {
	ListLatestPing []storage.Ping `json:"list_latest_ping"`
	Status         string         `json:"status"`
}

func NewListLatestPing(lister LatestPingLister) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		listLatestPing, err := lister.ListLatestPing(req.Context())
		if err != nil {
			err := fmt.Errorf("get list latest ping: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseListLatestPing{
			ListLatestPing: listLatestPing,
			Status:         "OK",
		}

		jsonResponseListLatestPing, err := json.Marshal(response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponseListLatestPing); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("List latest ping received successfully")
	}
}
