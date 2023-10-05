package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
)

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	alias := r.RequestURI
	log := slog.With(
		slog.String("request_id", r.Context().Value("requestID").(string)),
	)
	resURL, err := h.app.GetURL(r.Context(), alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			return
		}
		log.Error("failed to get url: " + err.Error())
	}

	log.Info("got url", slog.String("url", resURL))

	// redirect to found url
	http.Redirect(w, r, resURL, http.StatusFound)
}
