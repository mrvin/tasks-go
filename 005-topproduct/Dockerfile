## Build
FROM golang:1.19.2-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

WORKDIR /app

COPY go.mod ./
COPY Makefile ./

COPY *.go ./

RUN go build -o /topProduct

## Deploy
FROM alpine:3.16.2

WORKDIR /

COPY start.sh /start.sh
COPY testdata /testdata/
COPY --from=build /topProduct /topProduct

ENTRYPOINT ["/start.sh"]
