package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mrvin/tasks-go/notes/internal/storage"
)

func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	if _, err := s.insertUser.ExecContext(ctx, user.Name, user.HashPassword, user.Role); err != nil {
		//TODO: user already exists

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *Storage) GetUser(ctx context.Context, name string) (*storage.User, error) {
	var user storage.User

	if err := s.getUser.QueryRowContext(ctx, name).Scan(&user.HashPassword, &user.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", storage.ErrNoUser, name)
		}
		return nil, fmt.Errorf("can't scan user with name: %s: %w", name, err)
	}
	user.Name = name

	return &user, nil
}
