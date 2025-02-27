package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type WalletStorage interface {
	Create(ctx context.Context, balance float64) (uuid.UUID, error)
	Balance(ctx context.Context, walletID uuid.UUID) (float64, error)

	Send(ctx context.Context, transaction Transaction) error
	Deposit(ctx context.Context, walletID uuid.UUID, amount float64) error
	Withdraw(ctx context.Context, walletID uuid.UUID, amount float64) error

	HistoryTransactions(ctx context.Context, walletID uuid.UUID, limit, offset uint64) ([]Transaction, error)
}

type Transaction struct {
	Time         time.Time `json:"time"`
	WalletIDFrom uuid.UUID `json:"from"`
	WalletIDTo   uuid.UUID `json:"to"`
	Amount       float64   `json:"amount"`
}
