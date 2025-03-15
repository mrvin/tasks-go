package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
)

func (s *Storage) GetHistory(ctx context.Context, userName string) (storage.HistoryResponse, error) {
	slReceivedHistory := make([]storage.ReceivedTransactionResponse, 0)
	slSentHistory := make([]storage.SentTransactionResponse, 0)
	history := storage.HistoryResponse{slReceivedHistory, slSentHistory}

	rows, err := s.getReceivedHistory.QueryContext(ctx, userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return history, nil
		}
		return history, fmt.Errorf("can't get product: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var transaction storage.ReceivedTransactionResponse
		err = rows.Scan(&transaction.UserName, &transaction.Amount)
		if err != nil {
			return history, fmt.Errorf("can't scan next row: %w", err)
		}
		slReceivedHistory = append(slReceivedHistory, transaction)
	}
	if err := rows.Err(); err != nil {
		return history, fmt.Errorf("rows error: %w", err)
	}

	rows, err = s.getSentHistory.QueryContext(ctx, userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return history, nil
		}
		return history, fmt.Errorf("can't get product: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var transaction storage.SentTransactionResponse
		err = rows.Scan(&transaction.UserName, &transaction.Amount)
		if err != nil {
			return history, fmt.Errorf("can't scan next row: %w", err)
		}
		slSentHistory = append(slSentHistory, transaction)
	}
	if err := rows.Err(); err != nil {
		return history, fmt.Errorf("rows error: %w", err)
	}

	return storage.HistoryResponse{slReceivedHistory, slSentHistory}, nil
}
