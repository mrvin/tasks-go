package list

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mrvin/tasks-go/persons/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/persons/pkg/http/response"
)

const (
	defaultLimit  = 100
	defaultOffset = 0
)

type PersonLister interface {
	List(ctx context.Context, limit, offset uint64, ageFrom, ageTo uint64, gender, countryID string) ([]storage.Person, error)
}

type ResponsePersons struct {
	Persons []storage.Person `json:"persons"`
	Status  string           `json:"status"`
}

func New(lister PersonLister) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var err error
		limit := uint64(defaultLimit)
		limitStr := req.URL.Query().Get("limit")
		if limitStr != "" {
			limit, err = strconv.ParseUint(limitStr, 10, 64)
			if err != nil {
				err := fmt.Errorf("incorrect limit value: %w", err)
				slog.Error(err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
		}
		offset := uint64(defaultOffset)
		offsetStr := req.URL.Query().Get("offset")
		if offsetStr != "" {
			offset, err = strconv.ParseUint(offsetStr, 10, 64)
			if err != nil {
				err := fmt.Errorf("incorrect offset value: %w", err)
				slog.Error(err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
		}
		ageFrom := uint64(0)
		ageFromStr := req.URL.Query().Get("age_from")
		if ageFromStr != "" {
			ageFrom, err = strconv.ParseUint(ageFromStr, 10, 8)
			if err != nil {
				err := fmt.Errorf("incorrect age_from value: %w", err)
				slog.Error(err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
		}
		ageTo := uint64(150)
		ageToStr := req.URL.Query().Get("age_to")
		if ageToStr != "" {
			ageTo, err = strconv.ParseUint(ageToStr, 10, 8)
			if err != nil {
				err := fmt.Errorf("incorrect age_to value: %w", err)
				slog.Error(err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
		}
		gender := req.URL.Query().Get("gender")
		countryID := req.URL.Query().Get("country_id")

		persons, err := lister.List(req.Context(), limit, offset, ageFrom, ageTo, gender, countryID)
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
