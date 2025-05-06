package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mrvin/tasks-go/persons/internal/storage"
)

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
			return nil, fmt.Errorf("%w %v", storage.ErrNoPersonID, id)
		}
		return nil, fmt.Errorf("get person: %w", err)
	}

	return &person, nil
}
