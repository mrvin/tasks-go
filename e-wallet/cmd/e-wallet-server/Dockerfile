## Build
FROM golang:1.23.2-alpine AS build

LABEL maintainer="mrvin v.v.vinogradovv@gmail.com"

RUN apk update && apk add make && apk add tzdata

WORKDIR /app

# Copy the code into the container.
COPY cmd/e-wallet-server cmd/e-wallet-server
COPY internal internal
COPY pkg pkg

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

RUN cd cmd/e-wallet-server/ && make build

## Deploy
FROM scratch

WORKDIR /

COPY --from=build ["/usr/share/zoneinfo", "/usr/share/zoneinfo"]
COPY --from=build ["/app/bin/e-wallet-server", "/usr/local/bin/e-wallet-server"]

ENV TZ=Europe/Moscow

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/e-wallet-server"]
