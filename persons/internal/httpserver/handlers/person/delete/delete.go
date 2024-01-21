package delete

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	httpresponse "github.com/mrvin/tasks-go/persons/pkg/http/response"
)

type PersonDeleter interface {
	Delete(ctx context.Context, id int64) error
}

func New(deleter PersonDeleter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		idStr := req.URL.Query().Get("id")
		if idStr == "" {
			err := errors.New("id is empty")
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			err := fmt.Errorf("convert id: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := deleter.Delete(req.Context(), id); err != nil {
			err := fmt.Errorf("delete person from storage: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		httpresponse.WriteOK(res)

		slog.Info("Person deletion was successful")
	}
}
