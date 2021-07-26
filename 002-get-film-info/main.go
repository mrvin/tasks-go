// Problem from book 'The Go Programming Language. Alan A.A. Donovan,
// Brian W. Kernighan'. Exercis 4.13: The JSON-based web service of the
// Open Movie Database lets you search https://omdbapi.com/ for a movie by
// name and download its poster image. Write a tool poster that downloads
// the poster image for the movie named on the command line.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const requestTemplate = "http://www.omdbapi.com/?apikey=4g3v2195&t=Casablanca"

const (
	codeErrNotAPIKey = iota + 2
	codeErrNotFilmName
	codeErrRequestTempParse
	codeErrNotFoundFilm
	codeErrGetInfo
	codeSaveInfToFile
	codeErrGetPoster
)

var errNotFoundFilm = errors.New("not found film")

type infoFilm struct {
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

type options struct {
	apiKey     string
	filmName   string
	fullInf    bool
	year       int
	fileToSave string
}

func main() {
	var opt options

	flag.StringVar(&opt.apiKey, "k", "", "API key")
	flag.StringVar(&opt.filmName, "n", "", "movie title")
	flag.BoolVar(&opt.fullInf, "p", false, "full plot")
	flag.IntVar(&opt.year, "y", 0, "year of release")
	flag.StringVar(&opt.fileToSave, "f", "", "file path for saving information")

	flag.Parse()

	if opt.apiKey == "" {
		log.Printf("Key API not set. Flag '-k'.\n")
		os.Exit(codeErrNotAPIKey)
	}
	if opt.filmName == "" {
		log.Printf("Movie title not set. Flag '-n'.\n")
		os.Exit(codeErrNotFilmName)
	}

	query, err := queryBuild(&opt)
	if err != nil {
		log.Printf("Query build error: %v\n", err)
		os.Exit(codeErrRequestTempParse)
	}
	fmt.Printf("%s\n", query)

	inf, err := getInfo(query)
	if err != nil {
		log.Printf("Information get error: %v\n", err)
		if errors.Is(err, errNotFoundFilm) {
			os.Exit(codeErrNotFoundFilm)
		}
		os.Exit(codeErrGetInfo)
	}
	if opt.fileToSave != "" {
		if err := saveInfToFile(inf, opt.fileToSave); err != nil {
			log.Printf("File saving error: %v\n", err)
			os.Exit(codeSaveInfToFile)
		}
	}

	if inf.Poster != "N/A" {
		posterName := posterNameBuild(inf)
		size, err := getPoster(inf.Poster, posterName)
		if err != nil {
			log.Printf("Poster get error: %v\n", err)
			os.Exit(codeErrGetPoster)
		}
		fmt.Printf("Downloaded a file %s with size %d\n", posterName, size)
	}
	printInfFilm(inf)
}

func queryBuild(opt *options) (string, error) {
	u, err := url.Parse(requestTemplate)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("apikey", opt.apiKey)
	q.Set("t", opt.filmName)
	if opt.fullInf {
		q.Set("plot", "full")
	}
	if opt.year != 0 {
		q.Set("y", strconv.Itoa(opt.year))
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func getInfo(query string) (*infoFilm, error) {
	var inf infoFilm

	resp, err := http.Get(query) //nolint:gosec,bodyclose,noctx
	if resp != nil {
		defer closeHTTPResponse(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&inf); err != nil {
		return nil, err
	}

	if inf.Response == "False" {
		if inf.Error == "Movie not found!" {
			return nil, fmt.Errorf("%w: url:%s", errNotFoundFilm, query)
		}
		return nil, errors.New("response: false")
	}

	return &inf, nil
}

func saveInfToFile(inf *infoFilm, fileInfName string) error {
	var fileInf *os.File

	valueInf := reflect.ValueOf(*inf)
	typeOfInf := valueInf.Type()

	// TODO Add a check that the line already exists in the file. By
	//  calculating string hashes when opening a file.
	_, err := os.Stat(fileInfName)
	if os.IsNotExist(err) {
		if fileInf, err = os.Create(fileInfName); err != nil {
			return err
		}
		if _, err := fileInf.WriteString(buildTableHeader(typeOfInf)); err != nil {
			return err
		}
	} else if fileInf, err = os.OpenFile(fileInfName, os.O_APPEND|os.O_WRONLY, 0666); err != nil {
		return err
	}
	defer closeFile(fileInf)

	var str strings.Builder
	for i := 0; i < valueInf.NumField(); i++ {
		if _, ok := typeOfInf.Field(i).Tag.Lookup("notInput"); !ok {
			if valueInf.Field(i).Kind() == reflect.Slice {
				for j := 0; j < valueInf.Field(i).Len(); j++ {
					src := strings.ReplaceAll(valueInf.Field(i).Index(j).FieldByName("Source").String(),
						";", ".,")
					value := strings.ReplaceAll(valueInf.Field(i).Index(j).FieldByName("Value").String(),
						";", ".,")
					str.WriteString(src)
					str.WriteString(": ")
					str.WriteString(value)
					str.WriteRune('\t')
				}
			} else {
				str.WriteString(strings.ReplaceAll(valueInf.Field(i).String(), ";", ".,"))
			}
			str.WriteRune(';') // Column separator
		}
	}
	str.WriteRune('\n')

	if _, err := fileInf.WriteString(str.String()); err != nil {
		return err
	}

	if err := fileInf.Sync(); err != nil {
		return err
	}

	return nil
}

func buildTableHeader(typeOfInf reflect.Type) string {
	var tableHeader strings.Builder

	for i := 0; i < typeOfInf.NumField(); i++ {
		if _, ok := typeOfInf.Field(i).Tag.Lookup("notInput"); !ok {
			tableHeader.WriteString(typeOfInf.Field(i).Name)
			tableHeader.WriteRune(';')
		}
	}
	tableHeader.WriteRune('\n')

	return tableHeader.String()
}

func posterNameBuild(inf *infoFilm) string {
	var filmName strings.Builder

	filmName.WriteString("Poster_")
	filmName.WriteString(strings.ReplaceAll(inf.Title, " ", "_"))
	filmName.WriteRune('_')
	filmName.WriteString(inf.Year)
	filmName.WriteString(".jpg")

	return filmName.String()
}

func getPoster(url, filmName string) (int64, error) {
	filePoster, err := os.Create(filmName)
	if err != nil {
		return 0, err
	}
	defer closeFile(filePoster)

	resp, err := http.Get(url) //nolint:gosec,bodyclose,noctx
	if resp != nil {
		defer closeHTTPResponse(resp)
	}
	if err != nil {
		return 0, err
	}

	size, err := io.Copy(filePoster, resp.Body)
	if err != nil {
		return 0, err
	}

	return size, nil
}

func printInfFilm(inf *infoFilm) {
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

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func closeHTTPResponse(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}
