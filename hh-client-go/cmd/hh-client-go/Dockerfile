## Build
FROM golang:1.22.6-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make
RUN apk add --update ca-certificates

WORKDIR  /app

# Copy the code into the container.
COPY cmd/hh-client-go cmd/hh-client-go
COPY internal internal
COPY pkg pkg
COPY cert cert

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
#RUN go mod download

RUN cd cmd/hh-client-go/ && make build

## Deploy
FROM scratch

#WORKDIR /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build ["/app/bin/hh-client-go", "/usr/local/bin/hh-client-go"]

ENTRYPOINT ["/usr/local/bin/hh-client-go"]
