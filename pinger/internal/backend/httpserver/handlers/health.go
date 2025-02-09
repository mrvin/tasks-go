package handlers

import (
	"net/http"

	httpresponse "github.com/mrvin/tasks-go/pinger/pkg/http/response"
)

func Health(res http.ResponseWriter, _ *http.Request) {
	httpresponse.WriteOK(res, http.StatusOK)
}
