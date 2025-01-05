package mocks

import (
	"context"

	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
	"github.com/stretchr/testify/mock"
)

type WalletSender struct {
	mock.Mock
}

func NewWalletSender() *WalletSender {
	return new(WalletSender)
}

func (m *WalletSender) Send(ctx context.Context, transaction storage.Transaction) error {
	args := m.Called(ctx, transaction)

	if len(args) == 0 {
		panic("no return value specified for Send")
	}

	return args.Error(0) //nolint: wrapcheck
}
