FROM golang:latest AS builder

WORKDIR /app

COPY models models
COPY server server

WORKDIR /app/server

COPY server/go.mod server/go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

RUN go build -o main .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/server/main .

COPY --from=builder /app/server/docs /app/docs

EXPOSE 8000

ENTRYPOINT ["/app/main"]