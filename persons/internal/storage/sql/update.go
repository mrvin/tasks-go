package sqlstorage

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/mrvin/tasks-go/persons/internal/storage"
)

func (s *Storage) Update(ctx context.Context, id int64, person *storage.Person) error {
	// Create sql query
	sqlUpdate := squirrel.Update("persons")
	sqlUpdate = sqlUpdate.Where(squirrel.Eq{"id": id})
	if person.Name != "" {
		sqlUpdate = sqlUpdate.Set("name", person.Name)
	}
	if person.Surname != "" {
		sqlUpdate = sqlUpdate.Set("surname", person.Surname)
	}
	if person.Patronymic != "" {
		sqlUpdate = sqlUpdate.Set("patronymic", person.Patronymic)
	}
	if person.Age != 0 {
		sqlUpdate = sqlUpdate.Set("age", person.Age)
	}
	if person.Gender != "" {
		sqlUpdate = sqlUpdate.Set("gender", person.Gender)
	}
	if person.CountryID != "" {
		sqlUpdate = sqlUpdate.Set("country_id", person.CountryID)
	}
	sqlUpdatePerson, args, err := sqlUpdate.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("unable to build UPDATE query: %w", err)
	}

	fmt.Printf("\n\n%s\n\n", sqlUpdatePerson)

	// Exec sql query
	res, err := s.db.ExecContext(ctx, sqlUpdatePerson, args...)
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
