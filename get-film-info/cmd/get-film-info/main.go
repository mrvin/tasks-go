// Task from book 'The Go Programming Language. Alan A.A. Donovan,
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
	"log"
	"os"

	"github.com/mrvin/tasks-go/get-film-info/pkg/omdbapi"
)

const (
	codeErrNotAPIKey = iota + 2
	codeErrNotFilmName
	codeErrNotFoundFilm
	codeErrGetInfo
	codeSaveInfoToFile
	codeErrGetPoster
	codeErrJSONMarshaling
)

type options struct {
	apiKey        string
	filmTitle     string
	isFullPlot    bool
	yearOfRelease int
	fileToSave    string
}

func main() {
	var opt options

	flag.StringVar(&opt.apiKey, "k", "", "API key")
	flag.StringVar(&opt.filmTitle, "n", "", "movie title")
	flag.BoolVar(&opt.isFullPlot, "p", false, "full plot")
	flag.IntVar(&opt.yearOfRelease, "y", 0, "year of release")
	flag.StringVar(&opt.fileToSave, "f", "", "file path for saving information")

	flag.Parse()

	if opt.apiKey == "" {
		log.Printf("Key API not set. Flag '-k'.\n")
		os.Exit(codeErrNotAPIKey)
	}
	if opt.filmTitle == "" {
		log.Printf("Movie title not set. Flag '-n'.\n")
		os.Exit(codeErrNotFilmName)
	}

	omdbapiRequestURL := omdbapi.RequestBuild(opt.apiKey, opt.filmTitle, opt.isFullPlot, opt.yearOfRelease)

	fmt.Printf("OMDb API request URL: %s\n", omdbapiRequestURL)

	info, err := omdbapi.GetInfoFilm(omdbapiRequestURL)
	if err != nil {
		log.Printf("Error: get information: %v", err)
		if errors.Is(err, omdbapi.ErrNotFoundFilm) {
			os.Exit(codeErrNotFoundFilm)
		}
		os.Exit(codeErrGetInfo)
	}
	if opt.fileToSave != "" {
		if err := omdbapi.SaveInfoToFile(info, opt.fileToSave); err != nil {
			log.Printf("Error: file saving: %v", err)
			os.Exit(codeSaveInfoToFile)
		}
	}

	if info.Poster != "N/A" {
		posterName := omdbapi.PosterNameBuild(info)
		size, err := omdbapi.GetPoster(info.Poster, posterName)
		if err != nil {
			log.Printf("Error: get poster: %v", err)
			os.Exit(codeErrGetPoster)
		}
		fmt.Printf("Downloaded poster: %s with size %d\n", posterName, size)
	}

	jsonInfoFilm, err := json.MarshalIndent(info, "", "\t")
	if err != nil {
		log.Printf("Error: json marshaling failed: %v", err)
		os.Exit(codeErrJSONMarshaling)
	}
	fmt.Printf("%s\n", jsonInfoFilm)
}
