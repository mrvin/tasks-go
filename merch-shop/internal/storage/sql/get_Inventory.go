package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
)

func (s *Storage) GetInventory(ctx context.Context, userName string) ([]storage.ProductQuantity, error) {
	slProduct := make([]storage.ProductQuantity, 0)

	rows, err := s.getInventory.QueryContext(ctx, userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return slProduct, nil
		}
		return nil, fmt.Errorf("can't get product: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product storage.ProductQuantity
		err = rows.Scan(&product.Type, &product.Quantity)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		slProduct = append(slProduct, product)
	}
	if err := rows.Err(); err != nil {
		return slProduct, fmt.Errorf("rows error: %w", err)
	}

	return slProduct, nil
}
