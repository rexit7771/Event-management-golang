package structs

import "time"

type EventMessage struct {
	EventID   string      `json:"event_id"`
	Action    string      `json:"action"`
	EventData interface{} `json:"event_data"`
	UserID    uint        `json:"user_id"`
	Timestamp time.Time   `json:"timestamp"`
}
