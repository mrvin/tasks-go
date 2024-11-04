package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Add pure Go Postgres driver for the database/sql package.
	_ "github.com/lib/pq"
)

const maxOpenConns = 25
const maxIdleConns = 25
const connMaxLifetime = 5 // in minute

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

	insertBuilding *sql.Stmt
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
	s.db.SetConnMaxLifetime(connMaxLifetime * time.Minute)

	return nil
}

func (s *Storage) prepareQuery(ctx context.Context) error {
	var err error
	fmtStrErr := "prepare \"%s\" query: %w"

	sqlInsertBuilding := `
			INSERT INTO buildings (
				name,
				city,
				year,
				number_floors
			)
			VALUES ($1, $2, $3, $4)`
	s.insertBuilding, err = s.db.PrepareContext(ctx, sqlInsertBuilding)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert building", err)
	}

	return nil
}

func (s *Storage) Close() error {
	s.insertBuilding.Close()

	return s.db.Close() //nolint:wrapcheck
}
