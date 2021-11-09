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
	"html/template"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/mrvin/tasks-go/003-get-film-info-all/getinfofilm"
)

const requestTemplate = "http://www.omdbapi.com/?apikey=4g3v2195&t=Casablanca"

var imagesDir = "images/"

type options struct {
	apiKey            string
	fullInf           bool
	fileListFilmsName string
	fileSaveName      string
}

type shortInfoFilm struct {
	Title string
	Year  uint32
}

type slFilms []*getinfofilm.InfoFilm

func (films slFilms) Len() int           { return len(films) }
func (films slFilms) Less(i, j int) bool { return films[i].Title < films[j].Title }
func (films slFilms) Swap(i, j int)      { films[i], films[j] = films[j], films[i] }

type infoFilmsResult struct {
	Count uint64
	Films slFilms
}

var wgGetFilmInfo sync.WaitGroup

func main() {
	var opt options
	var listFilms []shortInfoFilm

	flag.StringVar(&opt.apiKey, "k", "", "API key")
	flag.BoolVar(&opt.fullInf, "p", false, "full plot")
	flag.StringVar(&opt.fileListFilmsName, "l", "", "file path list films")
	flag.StringVar(&opt.fileSaveName, "f", "", "file path for saving information")

	flag.Parse()

	if opt.apiKey == "" {
		log.Fatalf("Key API not set. Flag '-k'.\n")
	}
	if opt.fileListFilmsName == "" {
		log.Fatalf("File path list films not set. Flag '-l'.\n")
	}

	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		if err := os.Mkdir(imagesDir, 0755); err != nil {
			log.Fatalf("os.Mkdir: %v\n", err)
		}
	}
	if _, err := os.Stat("html"); os.IsNotExist(err) {
		if err := os.Mkdir("html", 0755); err != nil {
			log.Fatalf("os.Mkdir: %v\n", err)
		}
	}

	dataListFilms, err := ioutil.ReadFile(opt.fileListFilmsName)
	if err != nil {
		log.Fatalf("Read list films file error: %v\n", err)
	}

	if err := json.Unmarshal(dataListFilms, &listFilms); err != nil {
		log.Fatalf("Error unmarshal file %s: %v\n", opt.fileListFilmsName, err)
	}

	chResult := make(chan *getinfofilm.InfoFilm, len(listFilms))
	infoResult := infoFilmsResult{0, make([]*getinfofilm.InfoFilm, 0, len(listFilms))}
	for _, film := range listFilms {
		query, err := queryBuild(&opt, &film)
		if err != nil {
			log.Fatalf("Query build error: %v\n", err)
		}
		log.Printf("%s\n", query)

		wgGetFilmInfo.Add(1)
		go getFilmInfo(&query, chResult)
	}

	wgGetFilmInfo.Wait()
	close(chResult)

	for film := range chResult {
		if opt.fileSaveName != "" {
			if err := film.SaveInfoToFile(&opt.fileSaveName); err != nil {
				log.Printf("SaveInfoToFile: %v\n", err)
			}
		}
		infoResult.Films = append(infoResult.Films, film)
		infoResult.Count++
	}

	sort.Sort(infoResult.Films)

	if err := generationHTMLfiles(infoResult); err != nil {
		log.Fatalf("generationHTMLfiles: %v\n", err)
	}
}

func queryBuild(opt *options, film *shortInfoFilm) (string, error) {
	u, err := url.Parse(requestTemplate)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("apikey", opt.apiKey)
	q.Set("t", film.Title)
	if opt.fullInf {
		q.Set("plot", "full")
	}

	q.Set("y", fmt.Sprintf("%d", film.Year))

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func getFilmInfo(url *string, chResult chan<- *getinfofilm.InfoFilm) {
	var inf getinfofilm.InfoFilm

	defer wgGetFilmInfo.Done()

	if err := inf.GetInfo(url); err != nil {
		if errors.Is(err, getinfofilm.ErrNotFoundFilm) {
			log.Printf("%v\n", err)
			return
		}
		log.Printf("Information get error: %v\n", err)
		return
	}

	chResult <- &inf

	if inf.Poster != "N/A" {
		size, posterName, err := inf.GetPoster(&imagesDir)
		if err != nil {
			log.Printf("Poster get error: %v\n", err)
			return
		}
		log.Printf("Downloaded a file %s with size %d\n", *posterName, size)
	}
}

func filehtmlNameBuild(inf *getinfofilm.InfoFilm) string {
	var filmName strings.Builder

	filmName.WriteString("html/")
	filmName.WriteString(strings.ReplaceAll(inf.Title, " ", "_"))
	filmName.WriteRune('_')
	filmName.WriteString(inf.Year)
	filmName.WriteString(".html")

	return filmName.String()
}

func generationHTMLfiles(infoResult infoFilmsResult) error {
	var listFilmsTemp = template.Must(template.New("listFilmsHTML").
		Funcs(template.FuncMap{"wrapReplaceAllSpace": wrapReplaceAllSpace}).
		Parse(listFilmsHTML))

	mainTableFile, err := os.Create("html/main.html")
	if err != nil {
		return err
	}
	if err := listFilmsTemp.Execute(mainTableFile, infoResult); err != nil {
		return err
	}
	closeFile(mainTableFile)

	var infoFilmTemp = template.Must(template.New("infoFilmHTML").
		Funcs(template.FuncMap{"wrapReplaceAllSpace": wrapReplaceAllSpace}).
		Parse(infoFilmHTML))
	for _, film := range infoResult.Films {
		infoFilmFile, err := os.Create(filehtmlNameBuild(film))
		if err != nil {
			log.Printf("os.Create: %v\n", err)
			continue
		}

		if err := infoFilmTemp.Execute(infoFilmFile, film); err != nil {
			log.Printf("Execute %v\n", err)
			continue
		}

		closeFile(infoFilmFile)
	}

	return nil
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func wrapReplaceAllSpace(str string) string {
	return strings.ReplaceAll(str, " ", "_")
}
