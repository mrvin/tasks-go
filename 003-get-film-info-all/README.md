## get-film-info

Problem from book 'The Go Programming Language. Alan A.A. Donovan,
Brian W. Kernighan'. Exercis 4.13: The JSON-based web service of the Open
Movie Database lets you search https://omdbapi.com/ for a movie by name
and download its poster image. Write a tool poster that downloads the
poster image for the movie named on the command line.

```shell script
$ make build 
go build -o get-film-inf-all
```
```shell script
$ ./get-film-inf-all -help
Usage of ./get-film-inf-all:
  -f string
    	file path for saving information
  -k string
    	API key
  -l string
    	file path list films
  -p	full plot

```
```shell script
$ ./get-film-inf-all -k 4g3v2195 -l "filmList.json" -p -f "result.json"
```
Example filmList.json file:
```json
[
	{
		"Title": "Gladiator",
		"Year": 2000
	},
	{
		"Title": "The Last Samurai",
		"Year": 2003
	}
]
```
