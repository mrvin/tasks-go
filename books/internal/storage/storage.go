package storage

import (
	"context"
)

type Book struct {
	Title   string
	Authors []string
}

type Storage interface {
	CreateBook(ctx context.Context, book *Book) error
	GetBookByTitle(ctx context.Context, title string) (*Book, error)
	ListBooksByAuthor(ctx context.Context, author string) ([]*Book, error)
	ListAllBooks(ctx context.Context) ([]*Book, error)
}
