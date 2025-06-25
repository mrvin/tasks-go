package reprioritize

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	natsmq "github.com/mrvin/tasks-go/goods/internal/queue/nats"
	"github.com/mrvin/tasks-go/goods/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/goods/pkg/http/response"
)

type Request struct {
	NewPriority int64 `json:"newPriority"`
}

type Response struct {
	Priorities []storage.Priority `json:"priorities"` // Все приоритеты, которые были изминены
}

type GoodReprioritizer interface {
	Reprioritize(ctx context.Context, id, projectID, newPriority int64) (*storage.Good, []storage.Priority, error)
}

func New(reprioritizer GoodReprioritizer, mq *natsmq.Queue) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "Reprioritize good: "

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

		// Read json request
		var request Request
		body, err := io.ReadAll(req.Body)
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		good, priorities, err := reprioritizer.Reprioritize(req.Context(), id, projectID, request.NewPriority)
		if err != nil {
			err := fmt.Errorf("reprioritize: %w", err)
			slog.Error(op + err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		response := Response{
			Priorities: priorities,
		}
		// Write json response
		jsonResponse, err := json.Marshal(&response)
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
			Description: "Reprioritize good",
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

		slog.Info("Reprioritize good", slog.Int64("id", good.ID))
	}
}
