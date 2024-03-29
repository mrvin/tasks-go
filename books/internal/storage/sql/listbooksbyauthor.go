package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/books/internal/storage"
)

func (s *Storage) ListBooksByAuthor(ctx context.Context, author string) ([]*storage.Book, error) {
	bookTitles := make([]string, 0)
	sqlGetBooksByAuthor := `
		SELECT b.title 
		  FROM book_author AS r
		  JOIN books AS b
		    ON r.id_book = b.id
		  JOIN authors AS a
		    ON r.id_author = a.id
		 WHERE a.name = ?`
	rows, err := s.db.QueryContext(
		ctx,
		sqlGetBooksByAuthor,
		author,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			listBooks := make([]*storage.Book, 0)
			return listBooks, nil
		}
		return nil, fmt.Errorf("can't get books: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var title string
		err = rows.Scan(&title)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		bookTitles = append(bookTitles, title)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	listBooks := make([]*storage.Book, 0)
	for _, title := range bookTitles {
		book, err := s.GetBookByTitle(ctx, title)
		if err != nil {
			return nil, fmt.Errorf("GetBookByTitle: %w")
		}
		listBooks = append(listBooks, book)
	}

	return listBooks, nil
}
