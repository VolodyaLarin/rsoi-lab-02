FROM golang:1.20 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go build -o ./.bin/bonus ./cmd/bonus/main.go
FROM debian:12
WORKDIR /root
RUN apt-get update && apt-get install ca-certificates -y && update-ca-certificates
COPY --from=build /app/.bin/bonus ./app
CMD ["./app"]
