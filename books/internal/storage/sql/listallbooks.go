package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/books/internal/storage"
)

func (s *Storage) ListAllBooks(ctx context.Context) ([]*storage.Book, error) {
	sqlGetAllBooks := `
		SELECT b.title, a.name
		  FROM book_author AS r
		  JOIN books AS b
		    ON r.id_book = b.id
		  JOIN authors AS a
		    ON r.id_author = a.id`
	rows, err := s.db.QueryContext(
		ctx,
		sqlGetAllBooks,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			listBooks := make([]*storage.Book, 0)
			return listBooks, nil
		}
		return nil, fmt.Errorf("can't get books: %w", err)
	}
	defer rows.Close()

	mBooks := make(map[string][]string)
	for rows.Next() {
		var bookTitle, bookAuthor string
		err = rows.Scan(&bookTitle, &bookAuthor)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		mBooks[bookTitle] = append(mBooks[bookTitle], bookAuthor)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	books := make([]*storage.Book, 0)
	for title, authors := range mBooks {
		books = append(books, &storage.Book{Title: title, Authors: authors})
	}

	return books, nil
}
