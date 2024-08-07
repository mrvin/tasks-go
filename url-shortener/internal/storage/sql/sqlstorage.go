package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
	// Add Go sqlite driver is a cgo-free for the database/sql package.
	_ "modernc.org/sqlite"
)

type Conf struct {
	Path string `yaml:"path"`
}

type Storage struct {
	db        *sql.DB
	insertURL *sql.Stmt
	getURL    *sql.Stmt
	deleteURL *sql.Stmt
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	var st Storage

	if err := st.connect(ctx, conf.Path); err != nil {
		return nil, err
	}
	if err := st.prepareQuery(ctx); err != nil {
		return nil, err
	}

	return &st, nil
}

func (s *Storage) connect(ctx context.Context, dbConfStr string) error {
	var err error
	s.db, err = sql.Open("sqlite", dbConfStr)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	return nil
}

func (s *Storage) prepareQuery(ctx context.Context) error {
	var err error
	fmtStrErr := "prepare \"%s\" query: %w"

	const sqlInsertURL = "INSERT INTO url(url, alias) VALUES(?, ?) returning id"
	s.insertURL, err = s.db.PrepareContext(ctx, sqlInsertURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert url", err)
	}
	const sqlGetURL = "SELECT url FROM url WHERE alias = ?"
	s.getURL, err = s.db.PrepareContext(ctx, sqlGetURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select url", err)
	}
	const sqlDeleteURL = "DELETE FROM url WHERE alias = ?"
	s.deleteURL, err = s.db.PrepareContext(ctx, sqlDeleteURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "delete url", err)
	}

	return nil
}

func (s *Storage) CreateURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	var id int64
	if err := s.insertURL.QueryRowContext(ctx, urlToSave, alias).Scan(&id); err != nil {
		return 0, fmt.Errorf("put url: %w", err)
	}

	return id, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	var resURL string

	if err := s.getURL.QueryRowContext(ctx, alias).Scan(&resURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("execute statement: %w", err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	res, err := s.deleteURL.ExecContext(ctx, alias)
	if err != nil {
		return fmt.Errorf("delete url: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete url: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w: %q", storage.ErrURLAliasIsNotExists, alias)
	}
	return nil
}

func (s *Storage) Close() error {
	s.insertURL.Close()
	s.getURL.Close()
	s.deleteURL.Close()

	return s.db.Close() //nolint:wrapcheck
}
