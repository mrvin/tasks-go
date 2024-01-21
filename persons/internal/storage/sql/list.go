package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/persons/internal/storage"
)

func (s *Storage) List(ctx context.Context) ([]storage.Person, error) {
	persons := make([]storage.Person, 0)
	rows, err := s.listPersons.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return persons, nil
		}
		return nil, fmt.Errorf("can't get person: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var person storage.Person
		err = rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.CountryID,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		persons = append(persons, person)
	}
	if err := rows.Err(); err != nil {
		return persons, fmt.Errorf("rows error: %w", err)
	}

	return persons, nil
}
