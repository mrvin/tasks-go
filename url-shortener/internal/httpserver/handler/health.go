package handler

import (
	"net/http"
)

func (h *Handler) HealthLivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
