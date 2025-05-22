package list

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/quotes/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/quotes/pkg/http/response"
)

type QuoteLister interface {
	List(ctx context.Context, author string) ([]storage.Quote, error)
}

type ResponseQuotes struct {
	Quotes []storage.Quote `json:"quotes"`
	Status string          `json:"status"`
}

func New(lister QuoteLister) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "List quotes:"

		author := req.URL.Query().Get("author")

		quotes, err := lister.List(req.Context(), author)
		if err != nil {
			err := fmt.Errorf("get list quotes: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseQuotes{
			Quotes: quotes,
			Status: "OK",
		}
		jsonResponseQuotes, err := json.Marshal(response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponseQuotes); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("List quotes")
	}
}
