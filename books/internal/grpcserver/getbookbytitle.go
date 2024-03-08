package grpcserver

import (
	"context"
	"log/slog"

	"github.com/mrvin/tasks-go/books/internal/booksapi"
)

func (s *Server) GetBookByTitle(ctx context.Context, req *booksapi.Title) (*booksapi.Book, error) {

	book, err := s.st.GetBookByTitle(ctx, req.GetTitle())
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Search by title was successful")

	return &booksapi.Book{Title: book.Title, Authors: book.Authors}, nil
}
