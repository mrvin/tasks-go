## Build
FROM golang:1.23.3-alpine AS build

LABEL maintainer="mrvin v.v.vinogradovv@gmail.com"

RUN apk update && apk add make && apk add tzdata

WORKDIR /app

# Copy the code into the container.
COPY cmd/backend cmd/backend
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/backend/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=build ["/app/bin/backend", "/usr/local/bin/backend"]

ENV TZ=Europe/Moscow

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/backend"]
