package list

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/persons/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/persons/pkg/http/response"
)

type PersonLister interface {
	List(ctx context.Context) ([]storage.Person, error)
}

type ResponsePersons struct {
	Persons []storage.Person `json:"persons"`
	Status  string           `json:"status"`
}

func New(lister PersonLister) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		persons, err := lister.List(req.Context())
		if err != nil {
			err := fmt.Errorf("get list persons: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponsePersons{
			Persons: persons,
			Status:  "OK",
		}

		jsonResponsePersons, err := json.Marshal(response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponsePersons); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
