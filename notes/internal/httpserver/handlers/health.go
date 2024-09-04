package handlers

import (
	"net/http"

	"github.com/mrvin/tasks-go/notes/pkg/http/response"
)

func Health(res http.ResponseWriter, _ *http.Request) {
	response.WriteOK(res, http.StatusOK)
}
