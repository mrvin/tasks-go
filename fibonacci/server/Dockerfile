## Build
FROM golang:1.19.2-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

WORKDIR /app

# Download necessary Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY cache cache/
COPY config config/
COPY fibonacci fibonacci/
COPY servgrpc servgrpc/
COPY servhttp servhttp/

RUN CGO_ENABLED=0 go build -o /app/bin/server-fib -ldflags '-w'


## Deploy
FROM scratch

WORKDIR /

EXPOSE 8080
EXPOSE 55555

COPY config.yml ./
COPY --from=build ["/app/bin/server-fib", "/"]

ENTRYPOINT ["/server-fib"]
