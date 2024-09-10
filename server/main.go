package main

import (
	"flag"
	. "github.com/yankokirill/ImageProcessor/messaging"
	_ "github.com/yankokirill/ImageProcessor/server/docs"
	"github.com/yankokirill/ImageProcessor/server/http"
	. "github.com/yankokirill/ImageProcessor/storage"
	"log"
)

// @title Task Management API
// @version 1.0
// @description This is a sample server for managing tasks.
// @host localhost:8000
// @BasePath /
// @schemes commit
func main() {
	addr := flag.String("addr", ":8000", "address for commit server")
	s := NewRamStorage()
	b := NewProducerRMQ()
	server := http.NewServer(s, b)
	log.Printf("Starting server on %s", *addr)
	if err := http.CreateAndRunServer(server, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
