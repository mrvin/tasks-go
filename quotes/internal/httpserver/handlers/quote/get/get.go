package get

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/quotes/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/quotes/pkg/http/response"
)

type QuoteGetter interface {
	GetRandom(ctx context.Context) (*storage.Quote, error)
}

func New(getter QuoteGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "Get random quote: "

		quote, err := getter.GetRandom(req.Context())
		if err != nil {
			err := fmt.Errorf("get random quote: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		// Write json response
		jsonQuote, err := json.Marshal(&quote)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonQuote); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("Get random quote", slog.Int64("id", quote.ID))
	}
}
