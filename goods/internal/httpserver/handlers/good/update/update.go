package create

import (
	"context"
	"encoding/json"
	"errors"
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

type GoodUpdater interface {
	Update(ctx context.Context, id, projectID int64, name string, description string) (*storage.Good, error)
}

func New(updater GoodUpdater, mq *natsmq.Queue) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "Update good: "

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
		var request storage.Good
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

		good, err := updater.Update(req.Context(), id, projectID, request.Name, request.Description)
		if err != nil {
			err := fmt.Errorf("update good: %w", err)
			slog.Error(op + err.Error())
			if errors.Is(err, storage.ErrNoGoodID) {
				httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
				return
			}

			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		event := storage.Event{
			ID:          good.ID,
			ProjectID:   good.ProjectID,
			Name:        good.Name,
			Description: "Update good",
			Priority:    good.Priority,
			Removed:     good.Removed,
			Time:        time.Now(),
		}
		if err := mq.SendEvent(&event); err != nil {
			slog.Warn(op + "send event: " + err.Error())
		}

		// Write json response
		jsonResponse, err := json.Marshal(&good)
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

		slog.Info("Update good", slog.Int64("id", good.ID))
	}
}
