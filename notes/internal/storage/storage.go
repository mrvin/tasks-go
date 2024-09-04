package storage

import (
	"context"
	"errors"
)

var ErrNoUser = errors.New("no user with name")

//nolint:tagliatelle
type User struct {
	Name         string `json:"name"`
	HashPassword string `json:"hash_password"`
	Role         string `json:"role"`
}

type UserStorage interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, name string) (*User, error)
}

type Note struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type NotesStorage interface {
	CreateNote(ctx context.Context, userName string, note *Note) (int64, error)
	ListNotes(ctx context.Context, userName string) ([]Note, error)
}

type Storage interface {
	UserStorage
	NotesStorage
}
