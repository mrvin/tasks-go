package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
)

func (s *Storage) BuyProduct(ctx context.Context, userName, productName string) error {
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
	var balance uint64
	stmtGetBalance := tx.StmtContext(ctx, s.getBalance)
	defer stmtGetBalance.Close()
	if err := stmtGetBalance.QueryRowContext(ctx, userName).Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get balance: %w %v", storage.ErrNoFromUser, userName)
		}
		return fmt.Errorf("get balance: %w", err)
	}
	var price uint64
	stmtGetPrice := tx.StmtContext(ctx, s.getPrice)
	defer stmtGetPrice.Close()
	if err := stmtGetPrice.QueryRowContext(ctx, productName).Scan(&price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get product price: %w %v", storage.ErrProductNotFound, productName)
		}
		return fmt.Errorf("get product price: %w", err)
	}
	if balance < price {
		return storage.ErrNotEnoughFunds
	}

	stmtUpdateBalance := tx.StmtContext(ctx, s.updateBalance)
	defer stmtUpdateBalance.Close()
	if _, err := stmtUpdateBalance.ExecContext(ctx, userName, balance-price); err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	// Запись информации о заказе.
	stmtInsertOrder := tx.StmtContext(ctx, s.insertOrder)
	defer stmtInsertOrder.Close()
	if _, err := stmtInsertOrder.ExecContext(ctx, userName, productName); err != nil {
		return fmt.Errorf("recording order information: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
