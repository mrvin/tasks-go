package handler

import (
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/pkg/http/logger"
	"github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

func (h *Handler) DeleteURL(res http.ResponseWriter, req *http.Request) {
	log := slog.With(
		slog.String("Request ID", logger.GetRequestID(req.Context())),
	)
	alias := req.URL.Path[8:]
	if err := h.st.DeleteURL(req.Context(), alias); err != nil {
		log.Error("failed delete url: " + err.Error())
	}

	response.WriteOK(res, http.StatusOK)

	log.Info("url deleted", slog.String("url", alias))
}
