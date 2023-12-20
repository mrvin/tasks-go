package omdbapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var ErrNotFoundFilm = errors.New("not found film")

type InfoFilm struct {
	Title    string
	Year     string
	Rated    string
	Released string
	Runtime  string
	Genre    string
	Director string
	Writer   string
	Actors   string
	Plot     string
	Language string
	Country  string
	Awards   string
	Poster   string
	Ratings  []struct {
		Source string
		Value  string
	}
	Metascore  string
	ImdbRating string
	ImdbVotes  string
	ImdbID     string
	Type       string
	DVD        string
	BoxOffice  string
	Production string
	Website    string
	Response   string
	Error      string
}

// RequestBuild is builds a request to obtain information about a movie for the
// OMDb API.
func RequestBuild(apiKey, filmTitle string, isFullPlot bool, yearOfRelease int) string {
	requestURL := url.URL{
		Scheme: "https",
		Host:   "www.omdbapi.com",
	}

	v := url.Values{}
	v.Set("apikey", apiKey)
	v.Set("t", filmTitle)
	if isFullPlot {
		v.Set("plot", "full")
	}
	if yearOfRelease != 0 {
		v.Set("y", strconv.Itoa(yearOfRelease))
	}

	requestURL.RawQuery = v.Encode()

	return requestURL.String()
}

func GetInfoFilm(requestURL string) (*InfoFilm, error) {
	var info InfoFilm

	resp, err := http.Get(requestURL) //nolint:gosec,bodyclose,noctx
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, errors.New("unmarshal body response")
	}

	if info.Response == "False" {
		if info.Error == "Movie not found!" {
			return nil, fmt.Errorf("%w: url:%s", ErrNotFoundFilm, requestURL)
		}
		return nil, errors.New("response: false")
	}

	return &info, nil
}

func SaveInfoToFile(info *InfoFilm, fileInfoName string) error {
	data, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		return fmt.Errorf("JSON marshaling failed: %w", err)
	}

	fileInfo, err := os.OpenFile(fileInfoName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fileInfo.Close()

	fmt.Fprintf(fileInfo, "%s\n", data)

	return nil
}

func PosterNameBuild(info *InfoFilm) string {
	var filmName strings.Builder

	filmName.WriteString("./image/Poster_")
	filmName.WriteString(strings.ReplaceAll(info.Title, " ", "_"))
	filmName.WriteRune('_')
	filmName.WriteString(info.Year)
	filmName.WriteString(".jpg")

	return filmName.String()
}

func GetPoster(url, posterPath string) (int64, error) {
	resp, err := http.Get(url) //nolint:gosec,bodyclose,noctx
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(filepath.Dir(posterPath), 0770); err != nil {
		return 0, err
	}
	filePoster, err := os.Create(posterPath)
	if err != nil {
		return 0, err
	}
	defer filePoster.Close()

	size, err := io.Copy(filePoster, resp.Body)
	if err != nil {
		return 0, err
	}

	return size, nil
}
