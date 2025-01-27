package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
)

var ErrNoWalletIDFrom = errors.New("no wallet-from with id")
var ErrNoWalletIDTo = errors.New("no wallet-to with id")
var ErrNotEnoughFunds = errors.New("not enough funds in wallet")

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

	// Проверяем достаточно ли средств на исходящем кошельке
	var balanceFrom float64
	if err := tx.StmtContext(ctx, s.getBalance).QueryRow(transaction.WalletIDFrom).Scan(&balanceFrom); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get balance: %w %v", ErrNoWalletIDFrom, transaction.WalletIDFrom)
		}
		return fmt.Errorf("get balance: %w", err)
	}
	if balanceFrom-transaction.Amount < 0 {
		return ErrNotEnoughFunds
	}
	// Проверяем cуществует ли целевой кошелёк
	var balanceTo float64
	if err = tx.StmtContext(ctx, s.getBalance).QueryRow(transaction.WalletIDTo).Scan(&balanceTo); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get balance: %w %v", ErrNoWalletIDTo, transaction.WalletIDTo)
		}
		return fmt.Errorf("get balance: %w", err)
	}

	// Обновляем исходящий и целевой кошельки
	if _, err := tx.StmtContext(ctx, s.withdraw).Exec(transaction.WalletIDFrom, transaction.Amount); err != nil {
		return fmt.Errorf("update balance: %w", err)
	}
	if _, err := tx.StmtContext(ctx, s.deposit).Exec(transaction.WalletIDTo, transaction.Amount); err != nil {
		return err
	}
	if _, err := tx.StmtContext(ctx, s.insertTransaction).Exec(transaction.WalletIDFrom, transaction.WalletIDTo, transaction.Amount); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) SendOld(ctx context.Context, transaction storage.Transaction) error {
	if _, err := s.sendOld.ExecContext(
		ctx,
		transaction.WalletIDFrom,
		transaction.WalletIDTo,
		transaction.Amount,
	); err != nil {
		switch {
		case strings.Contains(err.Error(), "not enough funds"):
			return fmt.Errorf("%w: %w", ErrNotEnoughFunds, err)
		case strings.Contains(err.Error(), transaction.WalletIDFrom.String()):
			return fmt.Errorf("%w: %w", ErrNoWalletIDFrom, err)

		case strings.Contains(err.Error(), transaction.WalletIDTo.String()):
			return fmt.Errorf("%w: %w", ErrNoWalletIDTo, err)
		default:
			return fmt.Errorf("call send: %w", err)
		}
	}

	return nil
}
