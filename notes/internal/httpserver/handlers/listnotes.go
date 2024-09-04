package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/notes/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/notes/pkg/http/response"
)

type NoteLister interface {
	ListNotes(ctx context.Context, userName string) ([]storage.Note, error)
}

type ResponseListNotes struct {
	Notes  []storage.Note `json:"notes"`
	Status string         `json:"status"`
}

func NewListNotes(lister NoteLister) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userName, _ := req.Context().Value("userName").(string)
		notes, err := lister.ListNotes(req.Context(), userName)
		if err != nil {
			err := fmt.Errorf("ListNotes: get list notes: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseListNotes{
			Notes:  notes,
			Status: "OK",
		}

		jsonResponseListNotes, err := json.Marshal(response)
		if err != nil {
			err := fmt.Errorf("ListNotes: marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponseListNotes); err != nil {
			err := fmt.Errorf("ListNotes: write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
