package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Add pure Go Postgres driver for the database/sql package.
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 5 * time.Minute
)

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

	insertSchedule *sql.Stmt
	getListID      *sql.Stmt
	getSchedule    *sql.Stmt
	getAllTaking   *sql.Stmt
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
	s.db, err = sql.Open("pgx", dbConfStr)
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

	sqlInsertSchedule := `
	INSERT INTO schedules (
		name_medicine,
		num_per_day,
		times,
		all_life,
		begin_date,
		end_date,
		user_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`
	s.insertSchedule, err = s.db.PrepareContext(ctx, sqlInsertSchedule)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert schedule", err)
	}

	sqlGetListID := `
		SELECT id
		FROM schedules
		WHERE user_id = $1`
	s.getListID, err = s.db.PrepareContext(ctx, sqlGetListID)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "get list id", err)
	}

	sqlGetSchedule := `
		SELECT id,
			name_medicine,
			num_per_day,
			times,
			all_life,
			begin_date,
			end_date,
			user_id
		FROM schedules
		WHERE user_id = $1 AND id = $2`
	s.getSchedule, err = s.db.PrepareContext(ctx, sqlGetSchedule)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "get schedule", err)
	}

	sqlGetAllTaking := `
		SELECT name_medicine, times
		FROM schedules
		WHERE user_id = $1 AND (all_life = true OR $2 BETWEEN begin_date AND end_date)`
	s.getAllTaking, err = s.db.PrepareContext(ctx, sqlGetAllTaking)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "get all taking", err)
	}

	return nil
}

func (s *Storage) Close() error {
	s.insertSchedule.Close()
	s.getListID.Close()
	s.getSchedule.Close()
	s.getAllTaking.Close()

	return s.db.Close() //nolint:wrapcheck
}
