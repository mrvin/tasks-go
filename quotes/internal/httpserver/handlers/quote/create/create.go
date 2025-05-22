package create

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/quotes/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/quotes/pkg/http/response"
)

type QuoteCreator interface {
	Create(ctx context.Context, quote *storage.QuoteWithoutID) (int64, error)
}

type ResponseCreate struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

func New(creator QuoteCreator) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "Create quote: "

		// Read json request
		var request storage.QuoteWithoutID
		body, err := io.ReadAll(req.Body)
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := creator.Create(req.Context(), &request)
		if err != nil {
			err := fmt.Errorf("save quote: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseCreate{
			ID:     id,
			Status: "OK",
		}
		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("Create new quote", slog.Int64("id", id))
	}
}
