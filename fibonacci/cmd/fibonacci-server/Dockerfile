## Build
FROM golang:1.21.1-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

WORKDIR /app

COPY cmd/fibonacci-server cmd/fibonacci-server
COPY internal internal

# Download necessary Go modules
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/fibonacci-server/ && make build

## Deploy
FROM scratch

WORKDIR /

EXPOSE 8080
EXPOSE 55555

COPY --from=build ["/app/bin/fibonacci-server", "/"]

ENTRYPOINT ["/fibonacci-server"]
