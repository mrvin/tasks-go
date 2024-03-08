package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/books/internal/storage"
)

var ErrBookByTitleNotFound = errors.New("book not found with title")

func (s *Storage) GetBookByTitle(ctx context.Context, title string) (*storage.Book, error) {
	var book storage.Book
	book.Title = title
	sqlGetAuthors := `
		SELECT a.name 
		  FROM book_author AS r
		  JOIN books AS b
		    ON r.id_book = b.id
		  JOIN authors AS a
		    ON r.id_author = a.id
		 WHERE b.title = ?`
	rows, err := s.db.QueryContext(
		ctx,
		sqlGetAuthors,
		title,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: title %q", ErrBookByTitleNotFound, title)
		}
		return nil, fmt.Errorf("can't get authors: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var authorName string
		err = rows.Scan(&authorName)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		book.Authors = append(book.Authors, authorName)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return &book, nil
}
