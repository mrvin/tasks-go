package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/notes/internal/spelling"
	"github.com/mrvin/tasks-go/notes/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/notes/pkg/http/response"
)

type NoteCreator interface {
	CreateNote(ctx context.Context, userName string, note *storage.Note) (int64, error)
}

type RequestSaveNote struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ResponseSaveNote struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

func NewSaveNote(creator NoteCreator) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Read json request
		var request RequestSaveNote

		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("SaveNote: read body request: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("SaveNote: unmarshal body request: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		ok, err := spelling.Check(req.Context(), request.Description)
		if err != nil {
			err := fmt.Errorf("SaveNote: spell check: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		if !ok {
			err := fmt.Errorf("Text %q not pass the spell check", request.Description)
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		//nolint:exhaustruct
		note := storage.Note{
			Title:       request.Title,
			Description: request.Description,
		}

		userName, _ := req.Context().Value("userName").(string)
		id, err := creator.CreateNote(req.Context(), userName, &note)
		if err != nil {
			err := fmt.Errorf("SaveNote: saving note to storage: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseSaveNote{
			ID:     id,
			Status: "OK",
		}

		jsonResponseSaveNote, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("SaveNote: marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponseSaveNote); err != nil {
			err := fmt.Errorf("SaveNote: write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("New note created successfully")
	}
}
