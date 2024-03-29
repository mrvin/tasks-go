package grpcserver

import (
	"context"
	"log/slog"

	"github.com/mrvin/tasks-go/books/internal/booksapi"
)

func (s *Server) ListBooksByAuthor(ctx context.Context, req *booksapi.Author) (*booksapi.ListBooks, error) {

	books, err := s.st.ListBooksByAuthor(ctx, req.Author)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	pbBooks := make([]*booksapi.Book, len(books))

	for i, book := range books {
		pbBooks[i] = &booksapi.Book{
			Title:   book.Title,
			Authors: book.Authors,
		}
	}

	slog.Info("Search by author was successful")

	return &booksapi.ListBooks{Books: pbBooks}, nil
}
