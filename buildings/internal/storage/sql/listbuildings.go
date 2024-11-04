package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/mrvin/tasks-go/buildings/internal/storage"
)

func (s *Storage) ListBuildings(ctx context.Context, city string, year int16, numFloors int16) ([]storage.Building, error) {
	buildings := make([]storage.Building, 0)
	sqBuildings := squirrel.Select("id", "name", "city", "year", "number_floors").From("buildings")
	if city != "" {
		sqBuildings = sqBuildings.Where(squirrel.Eq{"city": city})
	}
	if year != 0 {
		sqBuildings = sqBuildings.Where(squirrel.Eq{"year": year})
	}
	if numFloors != 0 {
		sqBuildings = sqBuildings.Where(squirrel.Eq{"number_floors": numFloors})
	}
	sqlSelectBuildings, args, err := sqBuildings.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return buildings, fmt.Errorf("unable to build SELECT query: %w", err)
	}
	log.Printf("Executing SQL: %s", sqlSelectBuildings)

	rows, err := s.db.QueryContext(ctx, sqlSelectBuildings, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return buildings, nil
		}
		return nil, fmt.Errorf("can't get list buildings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var building storage.Building
		err = rows.Scan(
			&building.ID,
			&building.Name,
			&building.City,
			&building.Year,
			&building.NumFloors,
		)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		buildings = append(buildings, building)
	}
	if err := rows.Err(); err != nil {
		return buildings, fmt.Errorf("rows error: %w", err)
	}

	return buildings, nil
}
