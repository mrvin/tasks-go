## Build
FROM golang:1.23.3-alpine AS build

LABEL maintainer="mrvin v.v.vinogradovv@gmail.com"

RUN apk update && apk add make && apk add tzdata

WORKDIR /app

# Copy the code into the container.
COPY cmd/pinger cmd/pinger
COPY internal internal

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/pinger/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=build ["/app/bin/pinger", "/usr/local/bin/pinger"]

ENV TZ=Europe/Moscow

ENTRYPOINT ["/usr/local/bin/pinger"]
