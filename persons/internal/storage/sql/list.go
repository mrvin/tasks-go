package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/mrvin/tasks-go/persons/internal/storage"
)

func (s *Storage) List(ctx context.Context, limit, offset uint64, ageFrom, ageTo uint64, gender, countryID string) ([]storage.Person, error) {
	sqlList := squirrel.Select("id", "name", "surname", "patronymic", "age", "gender", "country_id").From("persons")
	if ageFrom != 0 {
		sqlList = sqlList.Where(squirrel.GtOrEq{"age": ageFrom})
	}
	if ageTo != 150 {
		sqlList = sqlList.Where(squirrel.LtOrEq{"age": ageTo})
	}
	if gender != "" {
		sqlList = sqlList.Where(squirrel.Eq{"gender": gender})
	}
	if countryID != "" {
		sqlList = sqlList.Where(squirrel.Eq{"country_id": countryID})
	}
	sqlList = sqlList.OrderBy("id")
	sqlList = sqlList.Limit(limit)
	sqlList = sqlList.Offset(offset)
	sqlSelectList, args, err := sqlList.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to build SELECT query: %w", err)
	}
	fmt.Printf("\n\n%s\n\n", sqlSelectList)

	persons := make([]storage.Person, 0)
	rows, err := s.db.QueryContext(ctx, sqlSelectList, args...)
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
