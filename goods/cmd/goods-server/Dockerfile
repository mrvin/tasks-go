## Build
FROM golang:1.24.0-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

WORKDIR  /app

# Copy the code into the container.
COPY cmd/goods-server cmd/goods-server
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/goods-server/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/app/bin/goods-server", "/"]

EXPOSE 8080

ENTRYPOINT ["/goods-server"]
