package storage

import (
	"context"
	"errors"
	"time"
)

var ErrNoGoodID = errors.New("no good with id")

type GoodsStorage interface {
	Create(ctx context.Context, projectID int64, name, description string) (*Good, error)
	Update(ctx context.Context, id, projectID int64, name, description string) (*Good, error)
	Delete(ctx context.Context, id, projectID int64) (*Good, error)
	List(ctx context.Context, limit, offset uint64) ([]Good, error)
	Meta(ctx context.Context) (int64, int64, error)
	Reprioritize(ctx context.Context, id, projectID, newPriority int64) (*Good, []Priority, error)
}

type Good struct {
	ID          int64     `json:"id"`                    // Уникальный идентификатор товара
	ProjectID   int64     `json:"projectID"`             // Идентификатор проекта (кампании)
	Name        string    `json:"name"`                  // Название товара
	Description string    `json:"description,omitempty"` // Описание товара (может быть пустым)
	Priority    int64     `json:"priority"`              // Приоритет товара (начинается с 1)
	Removed     bool      `json:"removed"`               // Флаг удаления (true - удален)
	CreatedAt   time.Time `json:"createdAt"`             // Дата и время создания
}

type Event struct {
	ID          int64     `json:"id"`          // Уникальный идентификатор товара
	ProjectID   int64     `json:"projectID"`   // Идентификатор проекта (кампании)
	Name        string    `json:"name"`        // Название товара
	Description string    `json:"description"` // Описание события
	Priority    int64     `json:"priority"`    // Приоритет товара (начинается с 1)
	Removed     bool      `json:"removed"`     // Флаг удаления (true - удален)
	Time        time.Time `json:"time"`        // Дата и время события
}

type Priority struct {
	ID       int64 `json:"id"`       // Уникальный идентификатор товара
	Priority int64 `json:"priority"` // Приоритет товара (начинается с 1)
}
