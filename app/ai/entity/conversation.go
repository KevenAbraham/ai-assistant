package entity

import "time"

type Role string

const (
	RoleUser       Role = "user"
	RoleAssistant  Role = "assistant"
	RoleSystem     Role = "system"
	RoleToolUse    Role = "tool_use"
	RoleToolResult Role = "tool_result"
)

type Message struct {
	ID        string
	Role      Role
	Content   string
	CreatedAt time.Time
}

type Conversation struct {
	ID        string
	SessionID string
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}
