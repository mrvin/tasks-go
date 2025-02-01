package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
)

func (s *Storage) Send(ctx context.Context, transaction storage.Transaction) error {
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
	var balanceFrom float64
	stmtGetBalance := tx.StmtContext(ctx, s.getBalance)
	defer stmtGetBalance.Close()
	if err := stmtGetBalance.QueryRowContext(ctx, transaction.WalletIDFrom).Scan(&balanceFrom); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get balance: %w %v", ErrNoWalletIDFrom, transaction.WalletIDFrom)
		}
		return fmt.Errorf("get balance: %w", err)
	}
	if balanceFrom-transaction.Amount < 0 {
		return ErrNotEnoughFunds
	}

	// Проверяем cуществует ли целевой кошелёк.
	var balanceTo float64
	if err = stmtGetBalance.QueryRowContext(ctx, transaction.WalletIDTo).Scan(&balanceTo); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get balance: %w %v", ErrNoWalletIDTo, transaction.WalletIDTo)
		}
		return fmt.Errorf("get balance: %w", err)
	}

	// Обновляем исходящий и целевой кошельки.
	stmtWithdraw := tx.StmtContext(ctx, s.withdraw)
	defer stmtWithdraw.Close()
	if _, err := stmtWithdraw.ExecContext(ctx, transaction.WalletIDFrom, transaction.Amount); err != nil {
		return fmt.Errorf("update source balance: %w", err)
	}
	stmtDeposit := tx.StmtContext(ctx, s.deposit)
	defer stmtDeposit.Close()
	if _, err := stmtDeposit.ExecContext(ctx, transaction.WalletIDTo, transaction.Amount); err != nil {
		return fmt.Errorf("update destination balance: %w", err)
	}

	// Запись информации о транзакции.
	stmtInsertTransaction := tx.StmtContext(ctx, s.insertTransaction)
	defer stmtInsertTransaction.Close()
	if _, err := stmtInsertTransaction.ExecContext(ctx, transaction.WalletIDFrom, transaction.WalletIDTo, transaction.Amount); err != nil {
		return fmt.Errorf("recording transaction information: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
