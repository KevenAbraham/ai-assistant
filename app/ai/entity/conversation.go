package entity

import "time"

// Role represents who sent a message in a conversation.
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

// Message is a single turn in a conversation.
type Message struct {
	ID        string
	Role      Role
	Content   string
	CreatedAt time.Time
}

// Conversation groups messages belonging to the same session.
type Conversation struct {
	ID        string
	SessionID string
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}
