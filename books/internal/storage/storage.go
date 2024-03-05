package storage

import (
	"context"
)

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Book struct {
	ID      int64    `json:"id"`
	Title   string   `json:"title"`
	Authors []Author `json:"authors"`
}

type Storage interface {
	CreateBook(ctx context.Context, book *Book) error
	GetBooksByAuthor(ctx context.Context, author string) ([]string, error)
	GetAuthorsByTitle(ctx context.Context, title string) ([]string, error)
}
