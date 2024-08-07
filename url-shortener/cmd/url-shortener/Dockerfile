## Build
FROM golang:1.22.5-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

WORKDIR /app

# Copy the code into the container.
COPY cmd/url-shortener cmd/url-shortener
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Download necessary Go modules
RUN go mod download

RUN cd cmd/url-shortener/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/app/bin/url-shortener", "/usr/local/bin/url-shortener"]

EXPOSE 8081

ENTRYPOINT ["/usr/local/bin/url-shortener"]
