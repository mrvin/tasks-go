package storage

import (
	"context"
	"errors"
)

var (
	ErrEmptyStorage = errors.New("empty storage")
	ErrNoQuoteID    = errors.New("no quote with id")
)

type QuoteStorage interface {
	Create(ctx context.Context, quote *QuoteWithoutID) (int64, error)
	GetRandom(ctx context.Context) (*Quote, error)
	Delete(ctx context.Context, id int64) error

	List(ctx context.Context, author string) ([]Quote, error)
}

type QuoteWithoutID struct {
	Author string `json:"author"`
	Text   string `json:"quote"`
}

type Quote struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Text   string `json:"quote"`
}
