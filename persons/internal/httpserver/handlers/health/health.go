package health

import (
	"net/http"

	httpresponse "github.com/mrvin/tasks-go/persons/pkg/http/response"
)

func Health(res http.ResponseWriter, _ *http.Request) {
	httpresponse.WriteOK(res, http.StatusOK)
}
