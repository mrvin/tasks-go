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

type HostLister interface {
	ListHost(ctx context.Context) ([]storage.Host, error)
}

//nolint:tagliatelle
type ResponseListHost struct {
	ListHost []storage.Host `json:"list_host"`
	Status   string         `json:"status"`
}

func NewListHost(lister HostLister) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		listHost, err := lister.ListHost(req.Context())
		if err != nil {
			err := fmt.Errorf("get list host ip: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseListHost{
			ListHost: listHost,
			Status:   "OK",
		}

		jsonResponseListHost, err := json.Marshal(response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponseListHost); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("List host received successfully")
	}
}
