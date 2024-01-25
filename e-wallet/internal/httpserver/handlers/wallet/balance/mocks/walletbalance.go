package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type WalletBalance struct {
	mock.Mock
}

func NewWalletBalance() *WalletBalance {
	return &WalletBalance{}
}

func (m *WalletBalance) Balance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	args := m.Called(ctx, walletID)

	if len(args) == 0 {
		panic("no return value specified for Balance")
	}

	return args.Get(0).(float64), args.Error(1)
}
