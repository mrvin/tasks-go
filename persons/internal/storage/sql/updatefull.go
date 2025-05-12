package sqlstorage

import (
	"context"
	"fmt"

	"github.com/mrvin/tasks-go/persons/internal/storage"
)

func (s *Storage) UpdateFull(ctx context.Context, id int64, person *storage.Person) error {
	res, err := s.updatePerson.ExecContext(
		ctx,
		id,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.CountryID,
	)
	if err != nil {
		return fmt.Errorf("update person: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update person: %w", err)
	}
	if count != 1 {
		return fmt.Errorf("%w %v", storage.ErrNoPersonID, id)
	}

	return nil
}
