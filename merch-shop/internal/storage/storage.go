package storage

import (
	"context"
	"errors"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrAccountExists   = errors.New("account exists")

	ErrNoFromUser     = errors.New("no from user")
	ErrNoToUser       = errors.New("no to user")
	ErrNotEnoughFunds = errors.New("not enough funds in balance")

	ErrProductNotFound = errors.New("product not found")
)

type ShopStorage interface {
	CreateAccount(ctx context.Context, userName, hashPassword string, startingBalance uint64) error
	GetAccount(ctx context.Context, userName string) (*Account, error)
	GetBalance(ctx context.Context, userName string) (uint64, error)

	SendCoin(ctx context.Context, transaction *Transaction) error

	BuyProduct(ctx context.Context, userName, productName string) error
	GetInventory(ctx context.Context, userName string) ([]ProductQuantity, error)

	GetHistory(ctx context.Context, userName string) (HistoryResponse, error)
}

type Account struct {
	Name         string
	HashPassword string
	Balance      uint64
}

type ProductQuantity struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type Transaction struct {
	FromUser string
	ToUser   string
	Amount   uint64
}

type ReceivedTransactionResponse struct {
	UserName string `json:"fromUser"`
	Amount   uint64 `json:"amount"`
}

type SentTransactionResponse struct {
	UserName string `json:"toUser"`
	Amount   uint64 `json:"amount"`
}

type HistoryResponse struct {
	Received []ReceivedTransactionResponse `json:"received"`
	Sent     []SentTransactionResponse     `json:"sent"`
}
