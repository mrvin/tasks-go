package grpcserver

import (
	"context"
	"log/slog"

	"github.com/mrvin/tasks-go/books/internal/booksapi"
)

func (s *Server) GetAuthorsByTitle(ctx context.Context, req *booksapi.GetAuthorsByTitleRequest) (*booksapi.GetAuthorsByTitleResponse, error) {

	authors, err := s.st.GetAuthorsByTitle(ctx, req.GetTitle())
	if err != nil {
		return nil, err
	}

	slog.Info("Search by title was successful")

	return &booksapi.GetAuthorsByTitleResponse{Authors: authors}, nil
}
