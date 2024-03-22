## Build
FROM golang:1.21.6-alpine AS build

MAINTAINER mrvin <v.v.vinogradovv@gmail.com>

RUN apk add --update make

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

COPY --from=build ["/app/bin/e-wallet-server", "/"]

EXPOSE 8088

ENTRYPOINT ["/e-wallet-server"]
