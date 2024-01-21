package update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mrvin/tasks-go/persons/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/persons/pkg/http/response"
)

type PersonUpdater interface {
	Update(ctx context.Context, id int64, person *storage.Person) error
}

type RequestUpdate struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	CountryID  string `json:"countryID"`
}

func New(updater PersonUpdater) http.HandlerFunc {
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

		// Read json request
		var request RequestUpdate

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

		person := storage.Person{
			Name:       request.Name,
			Surname:    request.Surname,
			Patronymic: request.Patronymic,
			Age:        request.Age,
			Gender:     request.Gender,
			CountryID:  request.CountryID,
		}

		if err := updater.Update(req.Context(), id, &person); err != nil {
			err := fmt.Errorf("update person: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		httpresponse.WriteOK(res)

		slog.Info("Person information update was successful")
	}
}
