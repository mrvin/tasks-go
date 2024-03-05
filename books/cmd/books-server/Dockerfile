## Build
FROM golang:1.22.0-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

WORKDIR  /app

# Copy the code into the container.
COPY cmd/books-server cmd/books-server
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/books-server/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/app/bin/books-server", "/"]

EXPOSE 55555

ENTRYPOINT ["/books-server"]
