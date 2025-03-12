package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/app"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/logger"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/medication-scheduler/pkg/http/response"
)

type ScheduleGetter interface {
	GetSchedule(ctx context.Context, userID uuid.UUID, scheduleID int64) (*storage.Schedule, error)
}

func NewGetSchedule(getter ScheduleGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "Get schedule: "
		ctx := req.Context()
		userIDStr := req.URL.Query().Get("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			err := fmt.Errorf("parse user_id: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := logger.WithUserID(ctx, userID.String())
		scheduleIDStr := req.URL.Query().Get("schedule_id")
		scheduleID, err := strconv.ParseInt(scheduleIDStr, 10, 64)
		if err != nil {
			err := fmt.Errorf("convert schedule_id: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		ctx = logger.WithScheduleID(ctx, scheduleID)

		schedule, err := getter.GetSchedule(ctx, userID, scheduleID)
		if err != nil {
			err := fmt.Errorf("get schedule from db: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		schedule.Times = app.ConvertTimesToStrings(schedule.TimesInt64)
		schedule.Status = "OK"

		jsonResponse, err := json.Marshal(schedule)
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

		slog.InfoContext(ctx, "Get schedule")
	}
}
