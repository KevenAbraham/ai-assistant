package repository

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
	apprepository "github.com/KevenAbraham/ai-assistant/app/ai/repository"
	"github.com/KevenAbraham/ai-assistant/internal/database"
)

// conversationRepoPg implements apprepository.ConversationRepository using PostgreSQL.
type conversationRepoPg struct {
	db *database.DB
}

// NewConversationRepository creates a new PostgreSQL-backed ConversationRepository.
func NewConversationRepository(db *database.DB) apprepository.ConversationRepository {
	return &conversationRepoPg{db: db}
}

// newUUID generates a random UUID v4 using crypto/rand (no external dependencies).
func newUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand unavailable: %v", err))
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant RFC 4122
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (r *conversationRepoPg) Save(ctx context.Context, conv *entity.Conversation) error {
	conn := r.db.Conn()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if conv.ID == "" {
		conv.ID = newUUID()
	}

	now := time.Now().UTC()
	_, err = tx.Exec(ctx, `
		INSERT INTO conversations (id, session_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET updated_at = EXCLUDED.updated_at
	`, conv.ID, conv.SessionID, now, now)
	if err != nil {
		return fmt.Errorf("upsert conversation: %w", err)
	}

	for i := range conv.Messages {
		if conv.Messages[i].ID == "" {
			conv.Messages[i].ID = newUUID()
			conv.Messages[i].CreatedAt = now
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO messages (id, conversation_id, role, content, created_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO NOTHING
		`, conv.Messages[i].ID, conv.ID, conv.Messages[i].Role, conv.Messages[i].Content, conv.Messages[i].CreatedAt)
		if err != nil {
			return fmt.Errorf("insert message: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *conversationRepoPg) FindByID(ctx context.Context, id string) (*entity.Conversation, error) {
	conn := r.db.Conn()
	var conv entity.Conversation
	err := conn.QueryRow(ctx, `
		SELECT id, session_id, created_at, updated_at FROM conversations WHERE id = $1
	`, id).Scan(&conv.ID, &conv.SessionID, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", entity.ErrConversationNotFound, err)
	}
	if err := r.loadMessages(ctx, conn, &conv); err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *conversationRepoPg) FindBySessionID(ctx context.Context, sessionID string) (*entity.Conversation, error) {
	conn := r.db.Conn()
	var conv entity.Conversation
	err := conn.QueryRow(ctx, `
		SELECT id, session_id, created_at, updated_at
		FROM conversations WHERE session_id = $1
		ORDER BY updated_at DESC LIMIT 1
	`, sessionID).Scan(&conv.ID, &conv.SessionID, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", entity.ErrConversationNotFound, err)
	}
	if err := r.loadMessages(ctx, conn, &conv); err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *conversationRepoPg) FindRecent(ctx context.Context, limit int) ([]*entity.Conversation, error) {
	rows, err := r.db.Conn().Query(ctx, `
		SELECT id, session_id, created_at, updated_at
		FROM conversations ORDER BY updated_at DESC LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent conversations: %w", err)
	}
	defer rows.Close()

	var convs []*entity.Conversation
	for rows.Next() {
		var conv entity.Conversation
		if err := rows.Scan(&conv.ID, &conv.SessionID, &conv.CreatedAt, &conv.UpdatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, &conv)
	}
	return convs, rows.Err()
}

func (r *conversationRepoPg) AppendMessage(ctx context.Context, conversationID string, msg entity.Message) error {
	_, err := r.db.Conn().Exec(ctx, `
		INSERT INTO messages (id, conversation_id, role, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, msg.ID, conversationID, msg.Role, msg.Content, msg.CreatedAt)
	return err
}

func (r *conversationRepoPg) loadMessages(ctx context.Context, conn *pgx.Conn, conv *entity.Conversation) error {
	rows, err := conn.Query(ctx, `
		SELECT id, role, content, created_at FROM messages
		WHERE conversation_id = $1 ORDER BY created_at ASC
	`, conv.ID)
	if err != nil {
		return fmt.Errorf("load messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(&msg.ID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			return err
		}
		conv.Messages = append(conv.Messages, msg)
	}
	return rows.Err()
}
