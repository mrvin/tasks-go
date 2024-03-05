package grpcserver

import (
	"context"
	"log/slog"

	"github.com/mrvin/tasks-go/books/internal/booksapi"
)

func (s *Server) GetBooksByAuthor(ctx context.Context, req *booksapi.GetBooksByAuthorRequest) (*booksapi.GetBooksByAuthorResponse, error) {

	books, err := s.st.GetBooksByAuthor(ctx, req.Author)
	if err != nil {
		return nil, err
	}

	slog.Info("Search by author was successful")

	return &booksapi.GetBooksByAuthorResponse{Titles: books}, nil
}
