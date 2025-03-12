package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrvin/tasks-go/medication-scheduler/internal/app"
	"github.com/mrvin/tasks-go/medication-scheduler/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/medication-scheduler/pkg/http/response"
)

type ScheduleCreator interface {
	CreateSchedule(ctx context.Context, schedule *storage.Schedule) (int64, error)
}

type ResponseCreateSchedule struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

const (
	takeMedicineFrom = 8 * time.Hour  // Время начало приема лекарств 8:00
	takeMedicineTo   = 22 * time.Hour // Время окончания приема лекарств 22:00
)

func NewCreateSchedule(creator ScheduleCreator) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		// Read json request
		var request storage.Schedule
		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.ErrorContext(ctx, "Create schedule: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.ErrorContext(ctx, "Create schedule: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		request.TimesInt64 = app.GenerateTimeTaking(takeMedicineFrom, takeMedicineTo, request.NumPerDay)

		id, err := creator.CreateSchedule(ctx, &request)
		if err != nil {
			slog.ErrorContext(ctx, "Create schedule: "+err.Error())
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
			slog.ErrorContext(ctx, "Create schedule: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(ctx, "Create schedule: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.InfoContext(ctx, "Create new schedule",
			slog.Int64("id", id),
		)
	}
}
