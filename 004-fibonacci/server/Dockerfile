FROM golang:1.17.8-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY config.yml ./

COPY *.go ./
COPY cache cache/
COPY config config/
COPY fibonacci fibonacci/
COPY servgrpc servgrpc/
COPY servhttp servhttp/

RUN go build

EXPOSE 8080
EXPOSE 55555

ENTRYPOINT ["/app/server"]
