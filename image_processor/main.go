package main

import (
	"encoding/json"
	. "hw/image_processor/processor"
	. "hw/messaging"
	. "hw/models"
	. "hw/storage"
	"os"
)

func main() {
	postgresConnString := os.Getenv("POSTGRES_CONN_STRING")
	rabbitMQAddr := os.Getenv("RABBITMQ_ADDR")

	db := NewPostgresTaskRepo(postgresConnString)
	c := NewConsumerRMQ(rabbitMQAddr)
	msgs := c.Consume()

	for msg := range msgs {
		var task Task
		err := json.Unmarshal(msg.Body, &task)
		if err != nil {
			db.UpdateTaskStatus(task.ID, "failed", "Failed to read task")
			continue
		}
		status, result := Process(task)
		db.UpdateTaskStatus(task.ID, status, result)

		_ = msg.Ack(false)
	}
}
