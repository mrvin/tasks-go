package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/tasks-go/medication-scheduler/internal/app"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/logger"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/medication-scheduler/pkg/http/response"
)

type ScheduleSaver interface {
	SaveSchedule(ctx context.Context, schedule *storage.Schedule) (int64, error)
}

type ResponseCreateSchedule struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

func NewCreateSchedule(saver ScheduleSaver) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		op := "Create schedule: "
		ctx := req.Context()

		// Read json request
		var request storage.Schedule
		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		ctx = logger.WithUserID(ctx, request.UserID.String())

		// Validation
		if request.NumPerDay <= 0 || request.NumPerDay > 57 {
			err := errors.New("number per day should be in the range [1,57]")
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		if !request.AllLife {
			if time.Time(request.BeginDate).After(time.Time(request.EndDate)) {
				err := errors.New("begin date must be before end date")
				slog.ErrorContext(ctx, op+err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
		}

		request.Times = app.GenerateTimes(request.NumPerDay)

		id, err := saver.SaveSchedule(ctx, &request)
		if err != nil {
			err := fmt.Errorf("save schedule: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseCreateSchedule{
			ID:     id,
			Status: "OK",
		}
		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(ctx, op+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.InfoContext(ctx, "Create new schedule")
	}
}
