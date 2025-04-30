## Build
FROM golang:1.24.0-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make
RUN apk add --update ca-certificates

WORKDIR  /app

# Copy the code into the container.
COPY cmd/persons-server cmd/persons-server
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/persons-server/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build ["/app/bin/persons-server", "/"]

EXPOSE 8080

ENTRYPOINT ["/persons-server"]
