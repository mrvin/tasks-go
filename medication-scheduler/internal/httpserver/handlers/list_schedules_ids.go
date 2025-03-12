package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	httpresponse "github.com/mrvin/tasks-go/medication-scheduler/pkg/http/response"
)

type ScheduleLister interface {
	ListSchedulesIDs(ctx context.Context, userID uuid.UUID) ([]int64, error)
}

//nolint:tagliatelle
type ResponseListSchedules struct {
	ListID []int64 `json:"list_ids"`
	Status string  `json:"status"`
}

func NewListSchedulesIDs(lister ScheduleLister) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		userIDStr := req.URL.Query().Get("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			err := fmt.Errorf("parse user_id: %w", err)
			slog.ErrorContext(ctx, "List schedules: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		list, err := lister.ListSchedulesIDs(ctx, userID)
		if err != nil {
			err := fmt.Errorf("get list from db: %w", err)
			slog.ErrorContext(ctx, "List schedules: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
		}

		// Write json response
		response := ResponseListSchedules{
			ListID: list,
			Status: "OK",
		}
		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.ErrorContext(ctx, "List schedules: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(ctx, "List schedules: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.InfoContext(ctx, "List of schedule IDs retrieved successfully")
	}
}
