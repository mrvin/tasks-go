package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

func (h *Handler) PutURL(w http.ResponseWriter, r *http.Request) {
	var req Request

	log := slog.With(
		slog.String("request_id", r.Context().Value("requestID").(string)),
	)

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Error("Read body req: " + err.Error())
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		log.Error("Unmarshal body req: " + err.Error())
		return
	}

	//validator

	id, err := h.app.PutURL(r.Context(), req.URL, req.Alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			return
		}
		log.Error("failed to add url: " + err.Error())
		return
	}
	h.res.Add("GET "+req.Alias, h.Redirect)
	h.res.Add("DELETE "+req.Alias, h.DeleteURL)

	log.Info("url added", slog.Int64("id", id))
}
