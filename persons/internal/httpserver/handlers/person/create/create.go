package create

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/persons/internal/enrich"
	"github.com/mrvin/tasks-go/persons/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/persons/pkg/http/response"
)

type PersonCreator interface {
	Create(ctx context.Context, person *storage.Person) (int64, error)
}

type RequestCreate struct {
	Name       string `example:"Dmitriy"    json:"name"`
	Surname    string `example:"Ushakov"    json:"surname"`
	Patronymic string `example:"Vasilevich" json:"patronymic,omitempty"`
}

type ResponseCreate struct {
	ID     int64  `example:"1"  json:"id"`
	Status string `example:"OK" json:"status"`
}

// New сreates a person creation handler.
//
//	@Summary			Create person
//	@Description		Create new person
//	@Tags			persons
//	@Accept			json
//	@Produce			json
//	@Param			input body RequestCreate true "person data"
//	@Success			201  {object} ResponseCreate
//	@Failure			400  {object}  response.RequestError
//	@Failure			500  {object}  response.RequestError
//	@Router			/persons [post]
func New(creator PersonCreator) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "Create schedule: "

		// Read json request
		var request RequestCreate
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

		age, err := enrich.GetAge(request.Name)
		if err != nil {
			err := fmt.Errorf("get age: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		gender, err := enrich.GetGender(request.Name)
		if err != nil {
			err := fmt.Errorf("get gender: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		countryID, err := enrich.GetCountryID(request.Name)
		if err != nil {
			err := fmt.Errorf("get country ID: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		//nolint: exhaustruct
		person := storage.Person{
			Name:       request.Name,
			Surname:    request.Surname,
			Patronymic: request.Patronymic,
			Age:        age,
			Gender:     gender,
			CountryID:  countryID,
		}
		id, err := creator.Create(req.Context(), &person)
		if err != nil {
			err := fmt.Errorf("save person: %w", err)
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

		slog.Info("Create new person", slog.Int64("id", id))
	}
}
