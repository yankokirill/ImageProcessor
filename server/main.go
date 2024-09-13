package main

import (
	"flag"
	. "hw/messaging"
	_ "hw/server/docs"
	"hw/server/http"
	. "hw/storage"
	"log"
	"os"
)

// @title Task Management API
// @version 2.0
// @description This is a sample server for managing tasks.
// @host localhost:8000
// @BasePath /
// @schemes http
func main() {
	postgresConnString := os.Getenv("POSTGRES_CONN_STRING")
	redisAddr := os.Getenv("REDIS_ADDR")
	jwtSecret := os.Getenv("JWT_SECRET")
	rabbitMQAddr := os.Getenv("RABBITMQ_ADDR")

	addr := flag.String("addr", ":8000", "address for server")
	s := NewDatabaseStorage(postgresConnString, redisAddr, jwtSecret)
	b := NewProducerRMQ(rabbitMQAddr)
	server := http.NewServer(s, b)
	log.Printf("Starting server on %s", *addr)
	if err := http.CreateAndRunServer(server, *addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
