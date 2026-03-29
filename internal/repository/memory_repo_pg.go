package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
	apprepository "github.com/KevenAbraham/ai-assistant/app/ai/repository"
	"github.com/KevenAbraham/ai-assistant/internal/database"
)

type memoryRepoPg struct {
	db *database.DB
}

func NewMemoryRepository(db *database.DB) apprepository.MemoryRepository {
	return &memoryRepoPg{db: db}
}

func (r *memoryRepoPg) Save(ctx context.Context, mem *entity.Memory) error {
	now := time.Now().UTC()
	_, err := r.db.Conn().Exec(ctx, `
		INSERT INTO memories (key, value, created_at, updated_at)
		VALUES ($1, $2, $3, $3)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at
	`, mem.Key, mem.Value, now)
	return err
}

func (r *memoryRepoPg) FindByKey(ctx context.Context, key string) (*entity.Memory, error) {
	var mem entity.Memory
	err := r.db.Conn().QueryRow(ctx, `
		SELECT key, value, created_at, updated_at FROM memories WHERE key = $1
	`, key).Scan(&mem.Key, &mem.Value, &mem.CreatedAt, &mem.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", entity.ErrMemoryNotFound, err)
	}
	return &mem, nil
}

func (r *memoryRepoPg) FindAll(ctx context.Context) ([]*entity.Memory, error) {
	rows, err := r.db.Conn().Query(ctx, `SELECT key, value, created_at, updated_at FROM memories ORDER BY key`)
	if err != nil {
		return nil, err
	}
	return scanMemories(rows)
}

func (r *memoryRepoPg) Search(ctx context.Context, query string) ([]*entity.Memory, error) {
	rows, err := r.db.Conn().Query(ctx, `
		SELECT key, value, created_at, updated_at FROM memories
		WHERE key ILIKE $1 OR value ILIKE $1
		ORDER BY key
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	return scanMemories(rows)
}

func (r *memoryRepoPg) Delete(ctx context.Context, key string) error {
	_, err := r.db.Conn().Exec(ctx, `DELETE FROM memories WHERE key = $1`, key)
	return err
}

func scanMemories(rows pgx.Rows) ([]*entity.Memory, error) {
	defer rows.Close()
	var mems []*entity.Memory
	for rows.Next() {
		var mem entity.Memory
		if err := rows.Scan(&mem.Key, &mem.Value, &mem.CreatedAt, &mem.UpdatedAt); err != nil {
			return nil, err
		}
		mems = append(mems, &mem)
	}
	return mems, rows.Err()
}
