package getinfofilm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
)

var ErrNotFoundFilm = errors.New("not found film")

var muFileInf sync.Mutex

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
	Response   string `notInput:""`
	Error      string `notInput:""`
}

func (inf *InfoFilm) GetInfo(query *string) error {
	resp, err := http.Get(*query)
	if resp != nil {
		defer closeHTTPResponse(resp)
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&inf); err != nil {
		return err
	}

	if inf.Response == "False" {
		if inf.Error == "Movie not found!" {
			return fmt.Errorf("%w: url:%s", ErrNotFoundFilm, *query)
		}
		return errors.New("response: false")
	}

	return nil
}

func (inf *InfoFilm) SaveInfoToFile(fileInfName *string) error {
	var fileInf *os.File

	_, err := os.Stat(*fileInfName)
	if os.IsNotExist(err) {
		if fileInf, err = os.Create(*fileInfName); err != nil {
			return err
		}
	} else if fileInf, err = os.OpenFile(*fileInfName, os.O_APPEND|os.O_WRONLY, 0666); err != nil {
		return err
	}
	defer closeFile(fileInf)

	dataInfoFilm, err := json.Marshal(inf)
	if err != nil {
		return err
	}

	muFileInf.Lock()
	if _, err := fileInf.Write(dataInfoFilm); err != nil {
		return err
	}
	if err := fileInf.Sync(); err != nil {
		return err
	}
	muFileInf.Unlock()

	return nil
}

func (inf *InfoFilm) GetPoster(dirName *string) (int64, *string, error) {
	posterName := inf.posterNameBuild(dirName)

	filePoster, err := os.Create(posterName)
	if err != nil {
		return 0, nil, err
	}
	defer closeFile(filePoster)

	resp, err := http.Get(inf.Poster)
	if resp != nil {
		defer closeHTTPResponse(resp)
	}
	if err != nil {
		return 0, nil, err
	}

	size, err := io.Copy(filePoster, resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return size, &posterName, nil
}

func (inf *InfoFilm) posterNameBuild(dirName *string) string {
	var filmName strings.Builder

	filmName.WriteString(*dirName)
	filmName.WriteString("Poster_")
	filmName.WriteString(strings.ReplaceAll(inf.Title, " ", "_"))
	filmName.WriteRune('_')
	filmName.WriteString(inf.Year)
	filmName.WriteString(".jpg")

	return filmName.String()
}

func (inf *InfoFilm) PrintInfFilm() {
	valueInf := reflect.ValueOf(*inf)
	typeOfInf := valueInf.Type()

	for i := 0; i < valueInf.NumField(); i++ {
		if _, ok := typeOfInf.Field(i).Tag.Lookup("notInput"); !ok {
			if valueInf.Field(i).Kind() == reflect.Slice {
				if valueInf.Field(i).Len() > 0 {
					fmt.Println("Ratings:")
				}
				for j := 0; j < valueInf.Field(i).Len(); j++ {
					fmt.Printf("\t%s: ", valueInf.Field(i).Index(j).FieldByName("Source").String())
					fmt.Printf("%s\n", valueInf.Field(i).Index(j).FieldByName("Value").String())
				}
			} else {
				fmt.Printf("%s: %s\n", typeOfInf.Field(i).Name, valueInf.Field(i).String())
			}
		}
	}
}

func closeHTTPResponse(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
