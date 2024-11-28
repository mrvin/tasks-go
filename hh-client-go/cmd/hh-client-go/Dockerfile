## Build
FROM golang:1.23.3-alpine AS build

LABEL maintainer="mrvin v.v.vinogradovv@gmail.com"

RUN apk update && apk add make && apk add ca-certificates && apk add tzdata

WORKDIR  /app

# Copy the code into the container.
COPY cmd/hh-client-go cmd/hh-client-go
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
#RUN go mod download

RUN cd cmd/hh-client-go/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=build ["/etc/ssl/certs/ca-certificates.crt", "/etc/ssl/certs/ca-certificates.crt"]
COPY --from=build ["/app/bin/hh-client-go", "/usr/local/bin/hh-client-go"]

ENV TZ=Europe/Moscow

ENTRYPOINT ["/usr/local/bin/hh-client-go"]
