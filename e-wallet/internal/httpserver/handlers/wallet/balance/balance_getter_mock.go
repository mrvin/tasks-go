package balance

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type WalletBalanceMock struct {
	mock.Mock
}

func NewWalletBalance() *WalletBalanceMock {
	return new(WalletBalanceMock)
}

func (m *WalletBalanceMock) Balance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	args := m.Called(ctx, walletID)

	if len(args) == 0 {
		panic("no return value specified for Balance")
	}
	balance, ok := args.Get(0).(float64)
	if !ok {
		panic("return value of wrong type")
	}

	return balance, args.Error(1) //nolint: wrapcheck
}
