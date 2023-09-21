## Build
FROM golang:1.21.1-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

WORKDIR /app

COPY go.mod ./
COPY Makefile ./

COPY *.go ./

RUN go build -o /top-product

## Deploy
FROM alpine:3.16.2

WORKDIR /

COPY start.sh /start.sh
COPY testdata /testdata/
COPY --from=build /top-product /top-product

ENTRYPOINT ["/start.sh"]
