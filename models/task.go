package models

import (
	"github.com/google/uuid"
)

type Filter struct {
	Name       string         `json:"name"`
	Parameters map[string]any `json:"parameters"`
}

type ImageProcessorPayload struct {
	Filter Filter `json:"filter"`
	Image  string `json:"image"`
}

type Task struct {
	ID      uuid.UUID `json:"task_id"`
	UserID  uuid.UUID `json:"user_id"`
	Payload ImageProcessorPayload
	Status  string `json:"status"`
	Result  string `json:"result"`
}

func (t *Task) GetFloatParameter(name string) (float64, bool) {
	value, exists := t.Payload.Filter.Parameters[name]
	if !exists {
		return 0, false
	}
	answer, ok := value.(float64)
	return answer, ok
}
