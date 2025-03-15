package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
)

func (s *Storage) SendCoin(ctx context.Context, transaction *storage.Transaction) error {
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

	// Проверяем достаточно ли средств на исходящем кошельке.
	var balanceFrom uint64
	stmtGetBalance := tx.StmtContext(ctx, s.getBalance)
	defer stmtGetBalance.Close()
	if err := stmtGetBalance.QueryRowContext(ctx, transaction.FromUser).Scan(&balanceFrom); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get balance: %w %v", storage.ErrNoFromUser, transaction.FromUser)
		}
		return fmt.Errorf("get balance: %w", err)
	}
	if balanceFrom < transaction.Amount {
		return storage.ErrNotEnoughFunds
	}

	// Проверяем cуществует ли целевой кошелёк.
	var balanceTo uint64
	if err = stmtGetBalance.QueryRowContext(ctx, transaction.ToUser).Scan(&balanceTo); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get balance: %w %v", storage.ErrNoToUser, transaction.ToUser)
		}
		return fmt.Errorf("get balance: %w", err)
	}

	// Обновляем исходящий и целевой кошельки.
	stmtUpdateBalance := tx.StmtContext(ctx, s.updateBalance)
	defer stmtUpdateBalance.Close()
	if _, err := stmtUpdateBalance.ExecContext(ctx, transaction.FromUser, balanceFrom-transaction.Amount); err != nil {
		return fmt.Errorf("update source balance: %w", err)
	}
	if _, err := stmtUpdateBalance.ExecContext(ctx, transaction.ToUser, balanceTo+transaction.Amount); err != nil {
		return fmt.Errorf("update destination balance: %w", err)
	}

	// Запись информации о транзакции.
	stmtInsertTransaction := tx.StmtContext(ctx, s.insertTransaction)
	defer stmtInsertTransaction.Close()
	if _, err := stmtInsertTransaction.ExecContext(ctx, transaction.FromUser, transaction.ToUser, transaction.Amount); err != nil {
		return fmt.Errorf("recording transaction information: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
