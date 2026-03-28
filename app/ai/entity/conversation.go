package entity

import "time"

// Role represents who sent a message in a conversation.
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
	// RoleToolUse and RoleToolResult store multi-turn tool call context so that
	// the next request can reconstruct the full conversation including tool exchanges.
	RoleToolUse    Role = "tool_use"
	RoleToolResult Role = "tool_result"
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
