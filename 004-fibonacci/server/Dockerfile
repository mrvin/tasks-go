FROM golang:1.17.8-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./

RUN go build

ENTRYPOINT ["/app/004-fibonacci"]
