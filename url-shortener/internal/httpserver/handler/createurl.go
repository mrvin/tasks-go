package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
	"github.com/mrvin/tasks-go/url-shortener/pkg/http/logger"
	"github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

type Request struct {
	URL   string `json:"url"             validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

func (h *Handler) CreateURL(res http.ResponseWriter, req *http.Request) {
	log := slog.With(
		slog.String("Request ID", logger.GetRequestID(req.Context())),
	)
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.Error("Read body req: " + err.Error())
		return
	}

	var request Request
	if err := json.Unmarshal(body, &request); err != nil {
		log.Error("Unmarshal body req: " + err.Error())
		return
	}

	id, err := h.st.CreateURL(req.Context(), request.URL, request.Alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", request.URL))
			return
		}
		log.Error("failed to add url: " + err.Error())
		return
	}

	response.WriteOK(res, http.StatusCreated)

	log.Info("url added", slog.Int64("id", id))
}
