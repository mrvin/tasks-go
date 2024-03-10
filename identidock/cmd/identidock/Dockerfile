## Build
FROM golang:1.22.1-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

WORKDIR /app

# Copy the code into the container.
COPY cmd/identidock cmd/identidock
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Download necessary Go modules
RUN go mod download

RUN cd cmd/identidock/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/app/bin/identidock", "/"]

EXPOSE 8888

ENTRYPOINT ["/identidock"]
