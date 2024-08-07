package handler

import (
	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
)

type Handler struct {
	st storage.Storage
}

func New(st storage.Storage) *Handler {
	return &Handler{
		st: st,
	}
}
