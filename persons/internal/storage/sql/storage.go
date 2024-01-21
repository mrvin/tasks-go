package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/mrvin/tasks-go/persons/pkg/retry"
)

const retriesConnect = 5

type Conf struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Storage struct {
	db *sql.DB

	conf *Conf

	insertPerson *sql.Stmt
	getPerson    *sql.Stmt
	updatePerson *sql.Stmt
	deletePerson *sql.Stmt
	listPersons  *sql.Stmt
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	var st Storage

	st.conf = conf

	if err := st.RetryConnect(ctx, retriesConnect); err != nil {
		return nil, fmt.Errorf("new database connection: %w", err)
	}

	if err := MigrationsUp(conf); err != nil {
		return nil, fmt.Errorf("database migrations: %w", err)
	}

	if err := st.prepareQuery(ctx); err != nil {
		return nil, fmt.Errorf("prepare query: %w", err)
	}

	return &st, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	dbConfStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.conf.Host, s.conf.Port, s.conf.User, s.conf.Password, s.conf.Name)
	s.db, err = sql.Open(s.conf.Driver, dbConfStr)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	// Setting db connections pool.
	s.db.SetMaxOpenConns(25)
	s.db.SetMaxIdleConns(25)
	s.db.SetConnMaxLifetime(5 * time.Minute)

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

	sqlInsertPerson := `
		INSERT INTO persons (
			name,
			surname,
			patronymic,
			age,
			gender,
			country_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	s.insertPerson, err = s.db.PrepareContext(ctx, sqlInsertPerson)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert person", err)
	}
	sqlGetPerson := `
		SELECT id, name, surname, patronymic, age, gender, country_id
		FROM persons
		WHERE id = $1`
	s.getPerson, err = s.db.PrepareContext(ctx, sqlGetPerson)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "get person", err)
	}
	sqlUpdatePerson := `
		UPDATE persons
		SET name = $2, surname = $3, patronymic = $4, age = $5, gender = $6, country_id = $7
		WHERE id = $1`
	s.updatePerson, err = s.db.PrepareContext(ctx, sqlUpdatePerson)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "update person", err)
	}
	sqlDeletePerson := `
		DELETE
		FROM persons
		WHERE id = $1`
	s.deletePerson, err = s.db.PrepareContext(ctx, sqlDeletePerson)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "delete person", err)
	}
	sqlListPersons := `
		SELECT id, name, surname, patronymic, age, gender, country_id
		FROM persons 
		ORDER BY id`
	s.listPersons, err = s.db.PrepareContext(ctx, sqlListPersons)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "list persons", err)
	}

	return nil
}

func (s *Storage) Close() error {
	s.insertPerson.Close()
	s.getPerson.Close()
	s.updatePerson.Close()
	s.deletePerson.Close()
	s.listPersons.Close()

	return s.db.Close() //nolint:wrapcheck
}
