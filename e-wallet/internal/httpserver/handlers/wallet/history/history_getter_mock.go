package history

import (
	"context"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
	"github.com/stretchr/testify/mock"
)

type WalletHistoryMock struct {
	mock.Mock
}

func NewWalletHistory() *WalletHistoryMock {
	return new(WalletHistoryMock)
}

func (m *WalletHistoryMock) HistoryTransactions(ctx context.Context, walletID uuid.UUID, limit, offset uint64) ([]storage.Transaction, error) {
	args := m.Called(ctx, walletID, limit, offset)

	if len(args) == 0 {
		panic("no return value specified for HistoryTransactions")
	}
	slTransactions, ok := args.Get(0).([]storage.Transaction)
	if !ok {
		panic("return value of wrong type")
	}

	return slTransactions, args.Error(1) //nolint: wrapcheck
}
