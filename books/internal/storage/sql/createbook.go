package sqlstorage

import (
	"context"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/mrvin/tasks-go/books/internal/storage"
)

func (s *Storage) CreateBook(ctx context.Context, book *storage.Book) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer tx.Rollback()

	sqlInsertBook := `
	INSERT INTO books (
		title
	) VALUES (?)`
	res, err := tx.ExecContext(
		ctx,
		sqlInsertBook,
		book.Title,
	)
	if err != nil {
		return fmt.Errorf("create book: %w", err)
	}

	bookID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id book: %w", err)
	}

	for _, author := range book.Authors {
		var authorID int64

		sqlInsertAuthor := `
		INSERT INTO authors (
			name
		) VALUES (?)`
		res, err := tx.ExecContext(
			ctx,
			sqlInsertAuthor,
			author,
		)
		if err != nil {
			if errMySQL, ok := err.(*mysql.MySQLError); ok {
				if errMySQL.Number == 1062 {
					authorID = s.getIDAuthor(ctx, author)
				} else {
					return fmt.Errorf("create author: %w", err)
				}
			}
		} else {
			authorID, err = res.LastInsertId()
			if err != nil {
				return fmt.Errorf("last insert id book: %w", err)
			}
		}

		sqlInsertBookAuthor := `
		INSERT INTO book_author (
			id_book,
			id_author
		) VALUES (?, ?)`
		if _, err := tx.ExecContext(
			ctx,
			sqlInsertBookAuthor,
			bookID,
			authorID,
		); err != nil {
			return fmt.Errorf("create book-author: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) getIDAuthor(ctx context.Context, name string) int64 {
	var authorID int64
	sqlGetIDAuthor := `SELECT id FROM authors WHERE name = ?`
	rows, _ := s.db.QueryContext(
		ctx,
		sqlGetIDAuthor,
		name,
	)
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&authorID)
	}

	return authorID
}
