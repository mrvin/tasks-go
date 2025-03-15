package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Add pure Go Postgres driver for the database/sql package.
	_ "github.com/lib/pq"
)

const (
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 5 * time.Minute
)

type Conf struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Storage struct {
	db *sql.DB

	conf *Conf

	insertAccount *sql.Stmt
	getAccount    *sql.Stmt
	getBalance    *sql.Stmt
	updateBalance *sql.Stmt

	insertOrder  *sql.Stmt
	getPrice     *sql.Stmt
	getInventory *sql.Stmt

	insertTransaction  *sql.Stmt
	getReceivedHistory *sql.Stmt
	getSentHistory     *sql.Stmt
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	var st Storage

	st.conf = conf

	if err := st.connect(ctx); err != nil {
		return nil, err
	}
	if err := st.prepareQuery(ctx); err != nil {
		return nil, err
	}

	return &st, nil
}

func (s *Storage) connect(ctx context.Context) error {
	var err error
	dbConfStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		s.conf.Host, s.conf.Port, s.conf.User, s.conf.Password, s.conf.Name)
	s.db, err = sql.Open("postgres", dbConfStr)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	// Setting db connections pool.
	s.db.SetMaxOpenConns(maxOpenConns)
	s.db.SetMaxIdleConns(maxIdleConns)
	s.db.SetConnMaxLifetime(connMaxLifetime)

	return nil
}

func (s *Storage) prepareQuery(ctx context.Context) error {
	var err error
	fmtStrErr := "prepare \"%s\" query: %w"

	// Account query.
	sqlInsertAccount := `
	INSERT INTO accounts (
		name,
		hash_password,
		balance
	) VALUES ($1, $2, $3)`
	s.insertAccount, err = s.db.PrepareContext(ctx, sqlInsertAccount)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert account", err)
	}
	sqlGetAccount := `
		SELECT name, hash_password, balance
		FROM accounts
		WHERE name = $1`
	s.getAccount, err = s.db.PrepareContext(ctx, sqlGetAccount)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "get account", err)
	}
	const sqlGetBalance = `
		SELECT balance
		FROM accounts
		WHERE name = $1`
	s.getBalance, err = s.db.PrepareContext(ctx, sqlGetBalance)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "get balance", err)
	}
	const sqlUpdateBalance = `
		UPDATE accounts
		SET balance = $2
		WHERE name = $1`
	s.updateBalance, err = s.db.PrepareContext(ctx, sqlUpdateBalance)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "update balance", err)
	}

	// Products query.
	const sqlInsertOrder = `
	INSERT INTO orders (
		user_name,
		product_name
	) VALUES ($1, $2)`
	s.insertOrder, err = s.db.PrepareContext(ctx, sqlInsertOrder)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert order", err)
	}
	const sqlGetPrice = `
		SELECT price
		FROM products
		WHERE name = $1`
	s.getPrice, err = s.db.PrepareContext(ctx, sqlGetPrice)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "get price", err)
	}
	const sqlGetInventory = `
		SELECT product_name, COUNT(*)
		FROM orders
		WHERE user_name = $1
		GROUP BY product_name`
	s.getInventory, err = s.db.PrepareContext(ctx, sqlGetInventory)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "get inventory", err)
	}

	// Transaction query.
	const sqlInsertTransaction = `
	INSERT INTO transactions (
		time,
		from_name,
		to_name,
		amount
	) VALUES (NOW(), $1, $2, $3)`
	s.insertTransaction, err = s.db.PrepareContext(ctx, sqlInsertTransaction)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insert transaction", err)
	}
	const sqlGetReceivedHistory = `
		SELECT from_name, amount
		FROM transactions
		WHERE to_name = $1`
	s.getReceivedHistory, err = s.db.PrepareContext(ctx, sqlGetReceivedHistory)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "get received history", err)
	}
	const sqlGetSentHistory = `
		SELECT to_name, amount
		FROM transactions
		WHERE from_name = $1`
	s.getSentHistory, err = s.db.PrepareContext(ctx, sqlGetSentHistory)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "get sent history", err)
	}

	return nil
}

func (s *Storage) Close() error {
	s.insertAccount.Close()
	s.getAccount.Close()
	s.getBalance.Close()
	s.updateBalance.Close()

	s.insertOrder.Close()
	s.getPrice.Close()
	s.getInventory.Close()

	s.insertTransaction.Close()
	s.getReceivedHistory.Close()
	s.getSentHistory.Close()

	return s.db.Close() //nolint:wrapcheck
}
