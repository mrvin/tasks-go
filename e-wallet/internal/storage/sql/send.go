package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
)

var ErrNoWalletIDFrom = errors.New("no wallet-from with id")
var ErrNoWalletIDTo = errors.New("no wallet-to with id")
var ErrNotEnoughFunds = errors.New("not enough funds in wallet")

func (s *Storage) Send(ctx context.Context, transaction storage.Transaction) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer tx.Rollback()

	// Проверяем достаточно ли средств на исходящем кошельке
	var balanceFrom float64
	sqlGetBalance := `
		SELECT balance
		FROM wallets
		WHERE id = $1`
	if err = tx.QueryRowContext(
		ctx,
		sqlGetBalance,
		transaction.WalletIDFrom,
	).Scan(&balanceFrom); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("get balance: %w %v", ErrNoWalletIDFrom, transaction.WalletIDFrom)
		}
		return fmt.Errorf("get balance: %w", err)
	}
	if balanceFrom-transaction.Amount < 0 {
		return ErrNotEnoughFunds
	}
	// Проверяем cуществует ли целевой кошелёк
	var balanceTo float64
	if err = tx.QueryRowContext(
		ctx,
		sqlGetBalance,
		transaction.WalletIDTo,
	).Scan(&balanceTo); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("get balance: %w %v", ErrNoWalletIDTo, transaction.WalletIDTo)
		}
		return fmt.Errorf("get balance: %w", err)
	}

	// Обновляем исходящий и целевой кошельки
	sqlUpdateBalance := `
		UPDATE wallets
		SET balance = $2
		WHERE id = $1`
	_, err = tx.ExecContext(
		ctx,
		sqlUpdateBalance,
		transaction.WalletIDFrom,
		balanceFrom-transaction.Amount,
	)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}
	_, err = tx.ExecContext(
		ctx,
		sqlUpdateBalance,
		transaction.WalletIDTo,
		balanceTo+transaction.Amount,
	)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}
	transaction.Time = time.Now()

	// Записываем транзакцию
	sqlInsertTransaction := `
		INSERT INTO transactions (
			time,
			from_wallet_id,
			to_wallet_id,
			amount
		)
		VALUES ($1, $2, $3, $4)`
	if _, err := tx.ExecContext(
		ctx,
		sqlInsertTransaction,
		transaction.Time,
		transaction.WalletIDFrom,
		transaction.WalletIDTo,
		transaction.Amount,
	); err != nil {
		return fmt.Errorf("write transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
