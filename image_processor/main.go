package main

import (
	"encoding/json"
	. "github.com/yankokirill/ImageProcessor/image_processor/commit"
	. "github.com/yankokirill/ImageProcessor/image_processor/processor"
	. "github.com/yankokirill/ImageProcessor/messaging"
	. "github.com/yankokirill/ImageProcessor/models"
)

func main() {
	c := NewConsumerRMQ()
	msgs := c.Consume()

	for msg := range msgs {
		var task Task
		err := json.Unmarshal(msg.Body, &task)
		if err != nil {
			CommitTask(task.ID, "failed", "Failed to read task")
			continue
		}
		status, result := Process(task)
		CommitTask(task.ID, status, result)

		_ = msg.Ack(false)
	}
}
