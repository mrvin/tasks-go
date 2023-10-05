package storage

import (
	"context"
	"errors"
)

var (
	ErrURLNotFound         = errors.New("url not found")
	ErrURLExists           = errors.New("url exists")
	ErrURLAliasIsNotExists = errors.New("url alias is not exists")
)

type Storage interface {
	PutURL(ctx context.Context, urlToSave string, alias string) (int64, error)
	GetURL(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, alias string) error
}
