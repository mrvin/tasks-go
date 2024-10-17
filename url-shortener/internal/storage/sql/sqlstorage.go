package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Add pure Go Postgres driver for the database/sql package.
	_ "github.com/lib/pq"
	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
)

const maxOpenConns = 25
const maxIdleConns = 25
const connMaxLifetime = 5 * time.Minute

type Conf struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Storage struct {
	db *sql.DB

	conf *Conf

	insertURL    *sql.Stmt
	selectGetURL *sql.Stmt
	deleteURL    *sql.Stmt
	getCount     *sql.Stmt
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	var st Storage

	st.conf = conf

	if err := st.connect(ctx); err != nil {
		return nil, err
	}
	if err := st.prepareQuery(ctx); err != nil {
		return nil, err
	}

	return &st, nil
}

func (s *Storage) connect(ctx context.Context) error {
	var err error
	dbConfStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		s.conf.Host, s.conf.Port, s.conf.User, s.conf.Password, s.conf.Name)
	s.db, err = sql.Open("postgres", dbConfStr)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	// Setting db connections pool.
	s.db.SetMaxOpenConns(maxOpenConns)
	s.db.SetMaxIdleConns(maxIdleConns)
	s.db.SetConnMaxLifetime(connMaxLifetime)

	return nil
}

func (s *Storage) prepareQuery(ctx context.Context) error {
	var err error
	fmtStrErr := "prepare \"%s\" query: %w"

	const sqlInsertURL = "INSERT INTO url(url, alias, count) VALUES($1, $2, 0) returning id"
	s.insertURL, err = s.db.PrepareContext(ctx, sqlInsertURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert url", err)
	}
	const sqlGetURL = "SELECT get_url($1)"
	s.selectGetURL, err = s.db.PrepareContext(ctx, sqlGetURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select get url", err)
	}
	const sqlDeleteURL = "DELETE FROM url WHERE alias = $1"
	s.deleteURL, err = s.db.PrepareContext(ctx, sqlDeleteURL)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "delete url", err)
	}
	const sqlGetCount = "SELECT count FROM url WHERE alias = $1"
	s.getCount, err = s.db.PrepareContext(ctx, sqlGetCount)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select count", err)
	}

	return nil
}

func (s *Storage) CreateURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	var id int64
	if err := s.insertURL.QueryRowContext(ctx, urlToSave, alias).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert url: %w", err)
	}

	return id, nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	var url string

	if err := s.selectGetURL.QueryRowContext(ctx, alias).Scan(&url); err != nil {
		return "", fmt.Errorf("execute statement: select get_url func: %w", err)
	}

	return url, nil
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

func (s *Storage) GetCount(ctx context.Context, alias string) (uint64, error) {
	var count uint64

	if err := s.getCount.QueryRowContext(ctx, alias).Scan(&count); err != nil {
		return 0, fmt.Errorf("execute statement: %w", err)
	}

	return count, nil
}

func (s *Storage) Close() error {
	s.insertURL.Close()
	s.selectGetURL.Close()
	s.deleteURL.Close()
	s.getCount.Close()

	return s.db.Close() //nolint:wrapcheck
}
