FROM golang:latest AS builder

WORKDIR /build

COPY models models
COPY storage storage
COPY messaging messaging
COPY image_processor image_processor


COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

WORKDIR /build/image_processor
RUN go build -o main .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /build/image_processor/main .

ENTRYPOINT ["/app/main"]