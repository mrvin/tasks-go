package create

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	natsmq "github.com/mrvin/tasks-go/goods/internal/queue/nats"
	"github.com/mrvin/tasks-go/goods/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/goods/pkg/http/response"
)

type DeleteResponse struct {
	ID        int64  `json:"id"`        // Уникальный идентификатор товара
	ProjectID int64  `json:"projectID"` // Идентификатор проекта (кампании)
	Removed   bool   `json:"removed"`   // Флаг удаления (true - удален)
	Status    string `json:"status"`
}

type GoodDeleter interface {
	Delete(ctx context.Context, id, projectID int64) (*storage.Good, error)
}

func New(deleter GoodDeleter, mq *natsmq.Queue) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "Delete good: "

		idStr := req.URL.Query().Get("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			err := fmt.Errorf("parse id: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		projectIDStr := req.URL.Query().Get("projectID")
		projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
		if err != nil {
			err := fmt.Errorf("parse projectID: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		good, err := deleter.Delete(req.Context(), id, projectID)
		if err != nil {
			err := fmt.Errorf("delete good: %w", err)
			slog.Error(op + err.Error())
			if errors.Is(err, storage.ErrNoGoodID) {
				httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
				return
			}

			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		deleteResponse := DeleteResponse{
			ID:        good.ID,
			ProjectID: good.ProjectID,
			Removed:   good.Removed,
			Status:    "OK",
		}
		jsonResponse, err := json.Marshal(&deleteResponse)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		event := storage.Event{
			ID:          good.ID,
			ProjectID:   good.ProjectID,
			Name:        good.Name,
			Description: "Delete good",
			Priority:    good.Priority,
			Removed:     good.Removed,
			Time:        time.Now(),
		}
		if err := mq.SendEvent(&event); err != nil {
			err := fmt.Errorf("send event: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("Delete good", slog.Int64("id", good.ID))
	}
}
