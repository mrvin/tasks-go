package handler

import (
	"net/http"

	"github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

func (h *Handler) Health(res http.ResponseWriter, _ *http.Request) {
	response.WriteOK(res, http.StatusOK)
}
