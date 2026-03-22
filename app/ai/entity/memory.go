package entity

import "time"

// Memory represents a long-term fact the assistant remembers about the user.
type Memory struct {
	ID        string
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
