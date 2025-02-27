package create

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/e-wallet/internal/app"
	httpresponse "github.com/mrvin/tasks-go/e-wallet/pkg/http/response"
)

type WalletCreator interface {
	Create(ctx context.Context, balance float64) (uuid.UUID, error)
}

type ResponseCreate struct {
	ID      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
	Status  string    `json:"status"`
}

func New(conf *app.Conf, creator WalletCreator) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id, err := creator.Create(req.Context(), conf.StartingBalance)
		if err != nil {
			err := fmt.Errorf("create wallet: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseCreate{
			ID:      id,
			Balance: conf.StartingBalance,
			Status:  "OK",
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("New wallet created successfully")
	}
}
