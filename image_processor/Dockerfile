FROM golang:latest AS builder

WORKDIR /build
COPY models models
COPY image_processor image_processor

WORKDIR /build/image_processor

COPY image_processor/go.mod image_processor/go.sum ./
RUN go mod download

RUN go build -o main .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /build/image_processor/main .

ENTRYPOINT ["/app/main"]