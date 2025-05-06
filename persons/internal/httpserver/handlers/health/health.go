package health

import (
	"net/http"

	httpresponse "github.com/mrvin/tasks-go/persons/pkg/http/response"
)

// Health is health check handler.
//
//	@Summary			Checking functionality
//	@Description		Checking functionality
//	@Tags			health
//	@Produce			json
//	@Success			200  {object} response.RequestOK
//	@Router			/health [get]
func Health(res http.ResponseWriter, _ *http.Request) {
	httpresponse.WriteOK(res, http.StatusOK)
}
