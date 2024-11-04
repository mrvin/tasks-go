package storage

import (
	"context"
)

//nolint:tagliatelle
type Building struct {
	ID        int64  `example:"0"                json:"id"`
	Name      string `example:"Building #1"      json:"name"`
	City      string `example:"Saint Petersburg" json:"city"`
	Year      int16  `example:"2022"             json:"year"`
	NumFloors int16  `example:"22"               json:"number_floors"`
}

type Storage interface {
	CreateBuilding(ctx context.Context, building *Building) error
	ListBuildings(ctx context.Context, city string, year int16, numFloors int16) ([]Building, error)
}
