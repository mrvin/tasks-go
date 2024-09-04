package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Add pure Go Postgres driver for the database/sql package.
	_ "github.com/lib/pq"
	"github.com/mrvin/tasks-go/notes/pkg/retry"
)

const retriesConnect = 5

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

	insertUser *sql.Stmt
	getUser    *sql.Stmt

	insertNote *sql.Stmt
	listNotes  *sql.Stmt
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	var st Storage

	st.conf = conf

	if err := st.RetryConnect(ctx, retriesConnect); err != nil {
		return nil, fmt.Errorf("new database connection: %w", err)
	}

	if err := st.prepareQuery(ctx); err != nil {
		return nil, fmt.Errorf("prepare query: %w", err)
	}

	return &st, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	dbConfStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		s.conf.Host, s.conf.Port, s.conf.User, s.conf.Password, s.conf.Name)
	s.db, err = sql.Open("postgres", dbConfStr)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	// Setting db connections pool.
	s.db.SetMaxOpenConns(maxOpenConns)
	s.db.SetMaxIdleConns(maxIdleConns)
	s.db.SetConnMaxLifetime(connMaxLifetime)

	return nil
}

func (s *Storage) RetryConnect(ctx context.Context, retries int) error {
	retryConnect := retry.Retry(s.Connect, retries)
	if err := retryConnect(ctx); err != nil {
		return fmt.Errorf("connection db: %w", err)
	}

	return nil
}

func (s *Storage) prepareQuery(ctx context.Context) error {
	var err error
	fmtStrErr := "prepare \"%s\" query: %w"

	// Users query.
	sqlInsertUser := `
		INSERT INTO users (
			name,
			hash_password,
			role
		)
		VALUES ($1, $2, $3)`
	s.insertUser, err = s.db.PrepareContext(ctx, sqlInsertUser)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insertUser", err)
	}
	sqlGetUser := `
		SELECT hash_password,
			role
		FROM users
		WHERE name = $1`
	s.getUser, err = s.db.PrepareContext(ctx, sqlGetUser)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "getUser", err)
	}

	// Notes query.
	sqlIsnsertNote := `
		INSERT INTO notes (
			title,
			description,
			user_name
		)
		VALUES ($1, $2, $3)
		RETURNING id`
	s.insertNote, err = s.db.PrepareContext(ctx, sqlIsnsertNote)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insertNote", err)
	}
	sqlListNotes := `
		SELECT id, title, description
		FROM notes 
		WHERE user_name = $1
		ORDER BY id`
	s.listNotes, err = s.db.PrepareContext(ctx, sqlListNotes)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "listNotes", err)
	}

	return nil
}

func (s *Storage) Close() error {
	s.insertUser.Close()
	s.getUser.Close()

	s.insertNote.Close()
	s.listNotes.Close()

	return s.db.Close() //nolint:wrapcheck
}
