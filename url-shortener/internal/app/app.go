package app

import (
	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
)

type App struct {
	storage.Storage
}

func New(st storage.Storage) *App {
	return &App{st}
}
