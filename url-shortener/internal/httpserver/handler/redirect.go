package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
	"github.com/mrvin/tasks-go/url-shortener/pkg/http/logger"
)

func (h *Handler) Redirect(res http.ResponseWriter, req *http.Request) {
	log := slog.With(
		slog.String("Request ID", logger.GetRequestID(req.Context())),
	)
	alias := req.URL.Path[8:]
	resURL, err := h.st.GetURL(req.Context(), alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			return
		}
		log.Error("failed to get url: " + err.Error())
	}

	log.Info("got url", slog.String("url", resURL))

	// redirect to found url
	http.Redirect(res, req, resURL, http.StatusFound)
}
