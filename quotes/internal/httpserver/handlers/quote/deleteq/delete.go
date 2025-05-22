package deletep

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	httpresponse "github.com/mrvin/tasks-go/quotes/pkg/http/response"
)

type QuoteDeleter interface {
	Delete(ctx context.Context, id int64) error
}

func New(deleter QuoteDeleter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "Delete quote: "

		idStr := req.PathValue("id")
		if idStr == "" {
			err := errors.New("id is empty")
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			err := fmt.Errorf("convert id: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := deleter.Delete(req.Context(), id); err != nil {
			err := fmt.Errorf("delete quote from storage: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		httpresponse.WriteOK(res, http.StatusOK)

		slog.Info("Delete quote", slog.Int64("id", id))
	}
}
