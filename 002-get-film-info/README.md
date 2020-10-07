## qwerty

This program gets information about the film from http://www.omdbapi.com/
and downloads the film poster. Film information can be written to a text file.
```shell script
$ make build 
go build -o get-film-inf
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
$ ./get-film-inf -k 4g3v2195 -n "Bullitt" -y 1968 -f "inf.txt"
http://www.omdbapi.com/?apikey=4g3v2195&t=Bullitt&y=1968
Downloaded a file Poster_Bullitt_1968.jpg with size 22404
Title: Bullitt
Year: 1968
...................
Production: Solar Productions
Website: N/A

```