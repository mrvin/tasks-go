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

	insertHost *sql.Stmt
	selectIP   *sql.Stmt

	insertPing       *sql.Stmt
	selectLatestPing *sql.Stmt
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

	// Host query.
	sqlInsertHost := `
		INSERT INTO hosts (
			name,
			ip
		)
		VALUES ($1, $2)`
	s.insertHost, err = s.db.PrepareContext(ctx, sqlInsertHost)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert host", err)
	}
	sqlSelectIP := `
		SELECT name, ip
		FROM hosts`
	s.selectIP, err = s.db.PrepareContext(ctx, sqlSelectIP)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select ip", err)
	}

	// Ping query.
	sqlInsertPing := `
		INSERT INTO pings (
			ip,
			time,
			created_at
		)
		VALUES ($1, $2, $3)`
	s.insertPing, err = s.db.PrepareContext(ctx, sqlInsertPing)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert ping", err)
	}
	sqlSelectLatestPing := `
		SELECT p.ip, p.time, mp.max_date
		FROM pings p
		JOIN(
			SELECT ip, MAX(created_at) AS max_date
			FROM pings
			WHERE time > 0
			GROUP BY ip
		) mp
		ON mp.ip = p.ip AND mp.max_date = p.created_at`
	s.selectLatestPing, err = s.db.PrepareContext(ctx, sqlSelectLatestPing)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "select ping", err)
	}

	return nil
}

func (s *Storage) Close() error {
	s.insertHost.Close()
	s.selectIP.Close()

	s.insertPing.Close()
	s.selectLatestPing.Close()

	return s.db.Close() //nolint:wrapcheck
}
