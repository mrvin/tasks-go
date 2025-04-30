package get

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mrvin/tasks-go/persons/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/persons/pkg/http/response"
)

type PersonGetter interface {
	Get(ctx context.Context, id int64) (*storage.Person, error)
}

func New(getter PersonGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
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

		person, err := getter.Get(req.Context(), id)
		if err != nil {
			err := fmt.Errorf("get person: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		jsonPerson, err := json.Marshal(&person)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonPerson); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("Person get successfully")
	}
}
