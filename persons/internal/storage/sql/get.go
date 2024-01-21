package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/persons/internal/storage"
)

var ErrNoPersonID = errors.New("no person with id")

func (s *Storage) Get(ctx context.Context, id int64) (*storage.Person, error) {
	var person storage.Person
	if err := s.getPerson.QueryRowContext(ctx, id).Scan(
		&person.ID,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.CountryID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("get person: %w %v", ErrNoPersonID, id)
		}
		return nil, fmt.Errorf("get person: %w", err)
	}

	return &person, nil
}
