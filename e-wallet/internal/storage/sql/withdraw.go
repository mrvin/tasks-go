package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (s *Storage) Withdraw(ctx context.Context, walletID uuid.UUID, amount float64) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if !errors.Is(err, sql.ErrTxDone) {
				slog.Error("Failed Rollback" + err.Error())
			}
		}
	}()

	// Проверяем достаточно ли средств на кошельке.
	var balanceFrom float64
	stmtGetBalance := tx.StmtContext(ctx, s.getBalance)
	defer stmtGetBalance.Close()
	if err := stmtGetBalance.QueryRowContext(ctx, walletID).Scan(&balanceFrom); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get balance: %w %v", ErrNoWalletID, walletID)
		}
		return fmt.Errorf("get balance: %w", err)
	}
	if balanceFrom-amount < 0 {
		return ErrNotEnoughFunds
	}

	// Обновляем кошелек.
	stmtWithdraw := tx.StmtContext(ctx, s.withdraw)
	defer stmtWithdraw.Close()
	if _, err := stmtWithdraw.ExecContext(ctx, walletID, amount); err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	// Запись информации о транзакции.
	stmtInsertTransaction := tx.StmtContext(ctx, s.insertTransaction)
	defer stmtInsertTransaction.Close()
	if _, err := stmtInsertTransaction.ExecContext(ctx, walletID, nil, amount); err != nil {
		return fmt.Errorf("recording transaction information: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
