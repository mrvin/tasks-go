## get-film-info

Task from book 'The Go Programming Language. Alan A.A. Donovan,
Brian W. Kernighan'. Exercis 4.13: The JSON-based web service of the Open
Movie Database lets you search https://omdbapi.com/ for a movie by name
and download its poster image. Write a tool poster that downloads the
poster image for the movie named on the command line.

```shell script
$ cd cmd/get-film-info/ && make build
go build -o ../../bin/get-film-info -ldflags '-w -s'
```
```shell script
$  ./bin/get-film-info -help
Usage of ./bin/get-film-info:
  -f string
    	file path for saving information
  -k string
    	API key
  -n string
    	movie title
  -p	full plot
  -y int
    	year of release
```
```shell script
$ ./bin/get-film-info -k 4g3v2195 -n "Bullitt" -y 1968 -f "movie_information.json"
OMDb API request URL: https://www.omdbapi.com?apikey=4g3v2195&t=Bullitt&y=1968
Downloaded poster: ./image/Poster_Bullitt_1968.jpg with size 25217
{
	"Title": "Bullitt",
	"Year": "1968",
	"Rated": "M/PG",
	"Released": "17 Oct 1968",
	"Runtime": "114 min",
	"Genre": "Action, Crime, Thriller",
	"Director": "Peter Yates",
	"Writer": "Alan Trustman, Harry Kleiner, Robert L. Fish",
	"Actors": "Steve McQueen, Jacqueline Bisset, Robert Vaughn",
	"Plot": "An all-guts, no-glory San Francisco cop becomes determined to find the underworld kingpin that killed the witness in his protection.",
	"Language": "English",
	"Country": "United States",
	"Awards": "Won 1 Oscar. 7 wins \u0026 9 nominations total",
	"Poster": "https://m.media-amazon.com/images/M/MV5BMGRjYzhmMGUtZjAyNy00NDkwLWI2ZmItMzJjYmYwM2JkYjkzXkEyXkFqcGdeQXVyMjUzOTY1NTc@._V1_SX300.jpg",
	"Ratings": [
		{
			"Source": "Internet Movie Database",
			"Value": "7.4/10"
		},
		{
			"Source": "Rotten Tomatoes",
			"Value": "98%"
		},
		{
			"Source": "Metacritic",
			"Value": "81/100"
		}
	],
	"Metascore": "81",
	"ImdbRating": "7.4",
	"ImdbVotes": "74,116",
	"ImdbID": "tt0062765",
	"Type": "movie",
	"DVD": "01 Sep 2008",
	"BoxOffice": "$511,350",
	"Production": "N/A",
	"Website": "N/A",
	"Response": "True",
	"Error": ""
}
```

```shell script
$ ./download_info_films.sh -k 4g3v2195 -f testdata/list_cartoons_films.txt
```