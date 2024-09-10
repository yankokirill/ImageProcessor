FROM golang:latest AS builder

WORKDIR /app

COPY models models
COPY messaging messaging
COPY storage storage
COPY server server

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

WORKDIR /app/server

RUN go build -o main .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/server/main .

EXPOSE 8000

ENTRYPOINT ["/app/main"]