package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
)

type Conf struct {
	Driver   string `yaml:"driver"`
	Path     string `yaml:"path"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Storage struct {
	db            *sql.DB
	stmtPutURL    *sql.Stmt
	stmtGetURL    *sql.Stmt
	stmtDeleteURL *sql.Stmt
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	var dbConfStr string
	switch conf.Driver {
	case "sqlite3":
		dbConfStr = conf.Path
	case "postgres":
		dbConfStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			conf.Host, conf.Port, conf.User, conf.Password, conf.Name)
	default:
		return nil, fmt.Errorf("driver does not support: %s", conf.Driver)
	}

	var st Storage
	if err := st.connect(ctx, conf.Driver, dbConfStr); err != nil {
		return nil, err
	}
	if err := st.createSchema(ctx); err != nil {
		return nil, err
	}
	if err := st.prepareQuery(ctx); err != nil {
		return nil, err
	}

	return &st, nil
}

func (s *Storage) connect(ctx context.Context, dbDriverStr string, dbConfStr string) error {
	var err error
	s.db, err = sql.Open(dbDriverStr, dbConfStr)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	return nil
}

func (s *Storage) createSchema(ctx context.Context) error {
	const createURLTable = `
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`
	stmt, err := s.db.Prepare(createURLTable)
	if err != nil {
		return fmt.Errorf("prepar create table url: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("exec create table url: %w", err)
	}

	return nil
}

func (s *Storage) prepareQuery(ctx context.Context) error {
	var err error
	fmtStrErr := "prepare \"%s\" query: %w"

	const sqlPutURL = "INSERT INTO url(url, alias) VALUES(?, ?) returning id"
	s.stmtPutURL, err = s.db.PrepareContext(ctx, sqlPutURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert url", err)
	}
	const sqlGetURL = "SELECT url FROM url WHERE alias = ?"
	s.stmtGetURL, err = s.db.PrepareContext(ctx, sqlGetURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select url", err)
	}
	const sqlDeleteURL = "DELETE FROM url WHERE alias = ?"
	s.stmtDeleteURL, err = s.db.PrepareContext(ctx, sqlDeleteURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "delete url", err)
	}

	return nil
}

func (s *Storage) PutURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	var id int64
	if err := s.stmtPutURL.QueryRowContext(ctx, urlToSave, alias).Scan(&id); err != nil {
		return 0, fmt.Errorf("put url: %w", err)
	}

	return id, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	var resURL string

	if err := s.stmtGetURL.QueryRowContext(ctx, alias).Scan(&resURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("Execute statement: %w", err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	res, err := s.stmtDeleteURL.ExecContext(ctx, alias)
	if err != nil {
		return fmt.Errorf("delete url: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete url: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w: %d", storage.ErrURLAliasIsNotExists, alias)
	}
	return nil
}

func (s *Storage) Close() error {
	s.stmtPutURL.Close()
	s.stmtGetURL.Close()
	s.stmtDeleteURL.Close()

	return s.db.Close() //nolint:wrapcheck
}
