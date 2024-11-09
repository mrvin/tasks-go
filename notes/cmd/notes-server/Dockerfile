## Build
FROM golang:1.23.3-alpine AS build

LABEL maintainer="mrvin v.v.vinogradovv@gmail.com"

RUN apk update && apk add make && apk add ca-certificates && apk add tzdata

WORKDIR  /app

# Copy the code into the container.
COPY cmd/notes-server cmd/notes-server
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/notes-server/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=build ["/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/ca-certificates.crt"]
COPY --from=build ["/app/bin/notes-server", "/usr/local/bin/notes-server"]

ENV TZ=Europe/Moscow

EXPOSE 8088

ENTRYPOINT ["/usr/local/bin/notes-server"]
