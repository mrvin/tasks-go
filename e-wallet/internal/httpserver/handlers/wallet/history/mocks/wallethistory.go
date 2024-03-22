package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
	"github.com/stretchr/testify/mock"
)

type WalletHistory struct {
	mock.Mock
}

func NewWalletHistory() *WalletHistory {
	return new(WalletHistory)
}

func (m *WalletHistory) HistoryTransactions(ctx context.Context, walletID uuid.UUID) ([]storage.Transaction, error) {
	args := m.Called(ctx, walletID)

	if len(args) == 0 {
		panic("no return value specified for HistoryTransactions")
	}

	return args.Get(0).([]storage.Transaction), args.Error(1)
}
