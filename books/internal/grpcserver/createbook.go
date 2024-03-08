package grpcserver

import (
	"context"
	"log/slog"

	"github.com/mrvin/tasks-go/books/internal/booksapi"
	"github.com/mrvin/tasks-go/books/internal/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CreateBook(ctx context.Context, req *booksapi.Book) (*emptypb.Empty, error) {

	book := storage.Book{
		Title: req.GetTitle(),
	}
	authors := req.GetAuthors()

	for _, author := range authors {
		book.Authors = append(book.Authors, author)
	}
	if err := s.st.CreateBook(ctx, &book); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("Save book was successful")

	return &emptypb.Empty{}, nil
}
