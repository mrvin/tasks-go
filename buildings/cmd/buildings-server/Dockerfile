## Build
FROM golang:1.23.2-alpine AS build

LABEL maintainer="mrvin v.v.vinogradovv@gmail.com"

RUN apk update && apk add make && apk add tzdata

WORKDIR /app

# Copy the code into the container.
COPY cmd/buildings-server cmd/buildings-server
COPY internal internal
COPY docs docs

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/buildings-server/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=build ["/app/bin/buildings-server", "/usr/local/bin/buildings-server"]

ENV TZ=Europe/Moscow

EXPOSE 8081

ENTRYPOINT ["/usr/local/bin/buildings-server"]
