package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

func (s *Storage) DepositOld(ctx context.Context, walletID uuid.UUID, amount float64) error {
	if _, err := s.depositOld.ExecContext(ctx, walletID, amount); err != nil {
		return fmt.Errorf("deposit: %w", err)
	}
	return nil
}

func (s *Storage) Deposit(ctx context.Context, walletID uuid.UUID, amount float64) error {
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

	if _, err := tx.StmtContext(ctx, s.deposit).Exec(walletID, amount); err != nil {
		return err
	}
	if _, err := tx.StmtContext(ctx, s.insertTransaction).Exec(nil, walletID, amount); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
