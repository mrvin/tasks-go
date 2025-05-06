package enrich

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseAge struct {
	Age int `json:"age"`
}

type ResponseGender struct {
	Gender string `json:"gender"`
}

type CountryIDandProbability struct {
	CountryID   string  `json:"country_id"` //nolint:tagliatelle
	Probability float32 `json:"probability"`
}

type ResponseCountryID struct {
	Country []CountryIDandProbability `json:"country"`
}

func GetAge(name string) (int, error) {
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

func GetGender(name string) (string, error) {
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

func GetCountryID(name string) (string, error) {
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
