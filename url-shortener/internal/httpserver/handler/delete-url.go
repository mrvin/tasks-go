package handler

import (
	"log/slog"
	"net/http"
)

func (h *Handler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	alias := r.RequestURI
	log := slog.With(
		slog.String("request_id", r.Context().Value("requestID").(string)),
	)
	if err := h.app.DeleteURL(r.Context(), alias); err != nil {
		log.Error("failed delete url: " + err.Error())
	}

	h.res.Delete("GET " + alias)
	h.res.Delete("DELETE " + alias)

	log.Info("url deleted", slog.String("url", alias))
}
