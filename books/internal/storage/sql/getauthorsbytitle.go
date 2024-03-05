package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) GetAuthorsByTitle(ctx context.Context, title string) ([]string, error) {
	authors := make([]string, 0)
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
			return authors, nil
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
		authors = append(authors, authorName)
	}
	if err := rows.Err(); err != nil {
		return authors, fmt.Errorf("rows error: %w", err)
	}

	return authors, nil
}
