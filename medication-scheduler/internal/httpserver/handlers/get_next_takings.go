package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/app"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/logger"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/medication-scheduler/pkg/http/response"
)

type AllTakingsGetter interface {
	GetAllTaking(ctx context.Context, userID uuid.UUID, now time.Time) ([]storage.AllTaking, error)
}

//nolint:tagliatelle
type ResponseGetNextTakings struct {
	TakingNow []storage.Taking `json:"taking_now"`
	Status    string           `json:"status"`
}

func NewGetNextTakings(getter AllTakingsGetter, periodNextTakings time.Duration) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "Get next takings: "
		ctx := req.Context()
		userIDStr := req.URL.Query().Get("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			err := fmt.Errorf("parse user_id: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		ctx = logger.WithUserID(ctx, userID.String())

		now := time.Now()
		allTaking, err := getter.GetAllTaking(ctx, userID, now)
		if err != nil {
			err := fmt.Errorf("get all taking from db: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseGetNextTakings{
			TakingNow: app.SelectNextTakings(allTaking, now, periodNextTakings),
			Status:    "OK",
		}
		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.InfoContext(ctx, "Get next takings medication")
	}
}
