package sqlstorage

import (
	"context"
	"fmt"

	"github.com/mrvin/tasks-go/persons/internal/storage"
)

func (s *Storage) Create(ctx context.Context, person *storage.Person) (int64, error) {
	var id int64
	if err := s.insertPerson.QueryRowContext(
		ctx,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.CountryID,
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("create: %w", err)
	}

	return id, nil
}
