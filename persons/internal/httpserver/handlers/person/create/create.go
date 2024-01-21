package create

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/persons/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/persons/pkg/http/response"
)

type PersonCreator interface {
	Create(ctx context.Context, person *storage.Person) (int64, error)
}

type RequestCreate struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type ResponseCreate struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

func New(creator PersonCreator) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Read json request
		var request RequestCreate

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

		age, err := getAge(request.Name)
		if err != nil {
			slog.Warn("get age: " + err.Error())
		}

		gender, err := getGender(request.Name)
		if err != nil {
			slog.Warn("get gender: " + err.Error())
		}

		countryID, err := getCountryID(request.Name)
		if err != nil {
			slog.Warn("get country ID: " + err.Error())
		}

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
			err := fmt.Errorf("create person: %w", err)
			slog.Error(err.Error())
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
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("New person created successfully")
	}
}

type ResponseAge struct {
	Age int `json:"age"`
}

func getAge(name string) (int, error) {
	var response ResponseAge
	const ageURL = "https://api.agify.io/?name="

	resp, err := http.Get(ageURL + name)
	if err != nil {
		return 0, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("http response status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("unmarshal body response: %w", err)
	}

	return response.Age, nil
}

type ResponseGender struct {
	Gender string `json:"gender"`
}

func getGender(name string) (string, error) {
	var response ResponseGender
	const genderURL = "https://api.genderize.io/?name="

	resp, err := http.Get(genderURL + name)
	if err != nil {
		return "", fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http response status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("unmarshal body response: %w", err)
	}

	return response.Gender, nil
}

type CountryIDandProbability struct {
	CountryID   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type ResponseCountryID struct {
	Country []CountryIDandProbability `json:"country"`
}

func getCountryID(name string) (string, error) {
	var response ResponseCountryID
	const countryIDURL = "https://api.nationalize.io/?name="

	resp, err := http.Get(countryIDURL + name)
	if err != nil {
		return "", fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http response status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("unmarshal body response: %w", err)
	}
	var maxProbability float32
	var maxIndex int
	for index, country := range response.Country {
		if country.Probability >= maxProbability {
			maxIndex = index
			maxProbability = country.Probability
		}
	}

	return response.Country[maxIndex].CountryID, nil
}
