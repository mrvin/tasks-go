package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
)

func (s *Storage) HistoryTransactions(ctx context.Context, walletID uuid.UUID) ([]storage.Transaction, error) {
	// Проверяем существует ли кошелек
	if err := s.getBalance.QueryRowContext(ctx, walletID).Scan(); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("check wallet: %w %v", ErrNoWalletID, walletID)
		}
	}

	transactions := make([]storage.Transaction, 0)
	rows, err := s.getHistory.QueryContext(ctx, walletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return transactions, nil
		}
		return nil, fmt.Errorf("can't get transactions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var transaction storage.Transaction
		err = rows.Scan(&transaction.Time, &transaction.WalletIDFrom, &transaction.WalletIDTo, &transaction.Amount)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return transactions, fmt.Errorf("rows error: %w", err)
	}

	return transactions, nil
}
