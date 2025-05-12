package updatefull

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

type PersonUpdaterFull interface {
	UpdateFull(ctx context.Context, id int64, person *storage.Person) error
}

type RequestUpdateFull struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	CountryID  string `json:"country_id"` //nolint:tagliatelle
}

// New сreates a person updation handler.
//
//	@Summary			Update person
//	@Description		Update all person information
//	@Tags			persons
//	@Accept			json
//	@Produce			json
//	@Param			id path int64 true "person id"
//	@Param			input body RequestUpdateFull true "person data"
//	@Success			200  {object} response.RequestOK
//	@Failure			400  {object}  response.RequestError
//	@Failure			500  {object}  response.RequestError
//	@Router			/persons/{id} [put]
func New(updater PersonUpdaterFull) http.HandlerFunc {
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

		// Read json request
		var request RequestUpdateFull
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

		//nolint: exhaustruct
		person := storage.Person{
			Name:       request.Name,
			Surname:    request.Surname,
			Patronymic: request.Patronymic,
			Age:        request.Age,
			Gender:     request.Gender,
			CountryID:  request.CountryID,
		}
		if err := updater.UpdateFull(req.Context(), id, &person); err != nil {
			err := fmt.Errorf("update person: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		httpresponse.WriteOK(res, http.StatusOK)

		slog.Info("Update full person", slog.Int64("id", id))
	}
}
