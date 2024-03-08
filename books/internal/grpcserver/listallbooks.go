package grpcserver

import (
	"context"
	"log/slog"

	"github.com/mrvin/tasks-go/books/internal/booksapi"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) ListAllBooks(ctx context.Context, _ *emptypb.Empty) (*booksapi.ListBooks, error) {
	books, err := s.st.ListAllBooks(ctx)
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

	return &booksapi.ListBooks{Books: pbBooks}, nil
}
