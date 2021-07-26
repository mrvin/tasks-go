## get-film-info

Problem from book 'The Go Programming Language. Alan A.A. Donovan,
Brian W. Kernighan'. Exercis 4.13: The JSON-based web service of the Open
Movie Database lets you search https://omdbapi.com/ for a movie by name
and download its poster image. Write a tool poster that downloads the
poster image for the movie named on the command line.

```shell script
$ make build 
go build -o get-film-inf
```
```shell script
$ ./get-film-inf -help
Usage of ./get-film-inf:
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
$ ./get-film-inf -k 4g3v2195 -n "Bullitt" -y 1968 -f "inf.txt"
http://www.omdbapi.com/?apikey=4g3v2195&t=Bullitt&y=1968
Downloaded a file Poster_Bullitt_1968.jpg with size 22404
Title: Bullitt
Year: 1968
...................
Production: Solar Productions
Website: N/A
```