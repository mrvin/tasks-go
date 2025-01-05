package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type WalletCreator struct {
	mock.Mock
}

func NewWalletCreator() *WalletCreator {
	return new(WalletCreator)
}

func (m *WalletCreator) Create(ctx context.Context, balance float64) (uuid.UUID, error) {
	args := m.Called(ctx, balance)

	if len(args) == 0 {
		panic("no return value specified for Create")
	}
	id, ok := args.Get(0).(uuid.UUID)
	if !ok {
		panic("return value of wrong type")
	}

	return id, args.Error(1) //nolint: wrapcheck
}
