package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	// Add pure Go Postgres driver for the database/sql package.
	_ "github.com/lib/pq"
	"github.com/mrvin/tasks-go/e-wallet/pkg/retry"
)

const retriesConnect = 5

const maxOpenConns = 25
const maxIdleConns = 25
const connMaxLifetime = 5 * time.Minute

var ErrNoWalletIDFrom = errors.New("no wallet-from with id")
var ErrNoWalletIDTo = errors.New("no wallet-to with id")
var ErrNoWalletID = errors.New("no wallet with id")
var ErrNotEnoughFunds = errors.New("not enough funds in wallet")

type Conf struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Storage struct {
	db *sql.DB

	conf *Conf

	insertWallet      *sql.Stmt
	insertTransaction *sql.Stmt
	getBalance        *sql.Stmt

	deposit  *sql.Stmt
	withdraw *sql.Stmt
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	var st Storage

	st.conf = conf

	if err := st.RetryConnect(ctx, retriesConnect); err != nil {
		return nil, fmt.Errorf("new database connection: %w", err)
	}

	if err := st.prepareQuery(ctx); err != nil {
		return nil, fmt.Errorf("prepare query: %w", err)
	}

	return &st, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	dbConfStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.conf.Host, s.conf.Port, s.conf.User, s.conf.Password, s.conf.Name)
	s.db, err = sql.Open(s.conf.Driver, dbConfStr)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	// Setting db connections pool.
	s.db.SetMaxOpenConns(maxOpenConns)
	s.db.SetMaxIdleConns(maxIdleConns)
	s.db.SetConnMaxLifetime(connMaxLifetime)

	return nil
}

func (s *Storage) RetryConnect(ctx context.Context, retries int) error {
	retryConnect := retry.Retry(s.Connect, retries)
	if err := retryConnect(ctx); err != nil {
		return fmt.Errorf("connection db: %w", err)
	}

	return nil
}

func (s *Storage) prepareQuery(ctx context.Context) error {
	var err error
	fmtStrErr := "prepare \"%s\" query: %w"

	const sqlInsertWallet = `
		INSERT INTO wallets (
			balance
		)
		VALUES ($1)
		RETURNING id`
	s.insertWallet, err = s.db.PrepareContext(ctx, sqlInsertWallet)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insertWallet", err)
	}
	const sqlGetBalance = `
		SELECT balance
		FROM wallets
		WHERE id = $1`
	s.getBalance, err = s.db.PrepareContext(ctx, sqlGetBalance)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "getBalance", err)
	}
	const sqlInsertTransaction = `
	INSERT INTO transactions (
		time,
		from_wallet_id,
		to_wallet_id,
		amount
	) VALUES (NOW(), $1, $2, $3)`
	s.insertTransaction, err = s.db.PrepareContext(ctx, sqlInsertTransaction)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insertTransaction", err)
	}

	const sqlWithdraw = `
		UPDATE wallets
		SET balance = round(CAST(balance-$2 AS numeric), 2)
		WHERE id = $1`
	s.withdraw, err = s.db.PrepareContext(ctx, sqlWithdraw)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "withdraw", err)
	}

	const sqlDeposit = `
		UPDATE wallets
		SET balance = round(CAST(balance+$2 AS numeric), 2)
		WHERE id = $1`
	s.deposit, err = s.db.PrepareContext(ctx, sqlDeposit)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "deposit", err)
	}

	return nil
}

func (s *Storage) Close() error {
	s.insertWallet.Close()
	s.insertTransaction.Close()
	s.getBalance.Close()

	s.withdraw.Close()
	s.deposit.Close()

	return s.db.Close() //nolint:wrapcheck
}
