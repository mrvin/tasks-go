build:
	go build -o get-film-inf-all
lint:
	golangci-lint run main.go listFilmsHTML.go infoFilmHTML.go || golangci-lint run ./getinfofilm/getinfofilm.go
run:
	go run .
fmt:
	go fmt && go fmt ./getinfofilm
import:
	goimports -w main.go && goimports -w ./getinfofilm/getinfofilm.go

