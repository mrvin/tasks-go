package storage

import (
	"context"
	"errors"
	"net"
	"time"
)

var (
	ErrHostExists = errors.New("host exists")
)

type PingerStorage interface {
	CreateHost(ctx context.Context, host *Host) error
	ListHost(ctx context.Context) ([]Host, error)

	CreatePing(ctx context.Context, ping *Ping) error
	ListLatestPing(ctx context.Context) ([]Ping, error)
}

type Host struct {
	Name string `json:"name"`
	IP   net.IP `json:"ip"`
}

//nolint:tagliatelle
type Ping struct {
	IP        net.IP        `json:"ip"`
	Time      time.Duration `json:"time"`
	CreatedAt time.Time     `json:"created_at"`
}
