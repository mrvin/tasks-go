package handler

import (
	"github.com/mrvin/tasks-go/url-shortener/internal/app"
	"github.com/mrvin/tasks-go/url-shortener/pkg/http/resolver"
)

type Handler struct {
	app *app.App
	res resolver.Resolver
}

func New(a *app.App, r resolver.Resolver) *Handler {
	return &Handler{
		app: a,
		res: r,
	}
}
