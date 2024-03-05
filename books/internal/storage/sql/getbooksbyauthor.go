package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) GetBooksByAuthor(ctx context.Context, author string) ([]string, error) {
	books := make([]string, 0)
	sqlGetBooks := `
		SELECT b.title 
		  FROM book_author AS r
		  JOIN books AS b
			ON r.id_book = b.id
		  JOIN authors AS a
			ON r.id_author = a.id
		 WHERE a.name = ?`
	rows, err := s.db.QueryContext(
		ctx,
		sqlGetBooks,
		author,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return books, nil
		}
		return nil, fmt.Errorf("can't get books: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var titleBook string
		err = rows.Scan(&titleBook)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		books = append(books, titleBook)
	}
	if err := rows.Err(); err != nil {
		return books, fmt.Errorf("rows error: %w", err)
	}

	return books, nil
}
