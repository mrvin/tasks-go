package sqlstorage

import (
	"context"
	"fmt"

	"github.com/mrvin/tasks-go/buildings/internal/storage"
)

func (s *Storage) CreateBuilding(ctx context.Context, building *storage.Building) error {
	if _, err := s.insertBuilding.ExecContext(ctx, building.Name, building.City, building.Year, building.NumFloors); err != nil {
		return fmt.Errorf("insert building: %w", err)
	}

	return nil
}
