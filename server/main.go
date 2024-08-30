package main

import (
	"flag"
	"github.com/streadway/amqp"
	"log"
	"server/http"
	"server/storage"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s", msg)
	}
}

// @title Task Management API
// @version 1.0
// @description This is a sample server for managing tasks.
// @host localhost:8000
// @BasePath /
// @schemes http
func main() {
	addr := flag.String("addr", ":8000", "address for http server")
	s := storage.NewRaiStorage()

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"task_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	server := http.NewServer(s, ch)

	log.Printf("Starting server on %s", *addr)
	if err := http.CreateAndRunServer(server, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
