package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	httpresponse "github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

type URLDeleter interface {
	DeleteURL(ctx context.Context, alias string) error
}

func NewDeleteURL(deleter URLDeleter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		alias := req.PathValue("alias")
		if err := deleter.DeleteURL(req.Context(), alias); err != nil {
			err := fmt.Errorf("failed delete url: %w", err)
			slog.ErrorContext(req.Context(), "Delete url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res, http.StatusOK)

		slog.InfoContext(req.Context(), "Deleted url", slog.String("alias", alias))
	}
}
