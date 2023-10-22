FROM golang:1.20 as build
WORKDIR /app
COPY go.mod go.sum ./
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go mod download
RUN go build -o ./.bin/ticket ./cmd/ticket/main.go
FROM debian:12-slim
WORKDIR /root
COPY --from=build /app/.bin/ticket ./app
CMD ["./app"]
