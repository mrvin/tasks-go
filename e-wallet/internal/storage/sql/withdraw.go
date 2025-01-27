package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (s *Storage) WithdrawOld(ctx context.Context, walletID uuid.UUID, amount float64) error {
	if _, err := s.withdrawOld.ExecContext(ctx, walletID, amount); err != nil {
		return fmt.Errorf("withdraw: %w", err)
	}
	return nil
}

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

	var balanceFrom float64
	if err := tx.StmtContext(ctx, s.getBalance).QueryRowContext(ctx, walletID).Scan(&balanceFrom); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get balance: %w %v", ErrNoWalletIDFrom, walletID)
		}
		return fmt.Errorf("get balance: %w", err)
	}
	if balanceFrom-amount < 0 {
		return ErrNotEnoughFunds
	}
	if _, err := tx.StmtContext(ctx, s.withdraw).Exec(walletID, amount); err != nil {
		return err
	}
	if _, err := tx.StmtContext(ctx, s.insertTransaction).Exec(walletID, nil, amount); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
