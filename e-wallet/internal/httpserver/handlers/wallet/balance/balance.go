package balance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	sqlstorage "github.com/mrvin/tasks-go/e-wallet/internal/storage/sql"
	httpresponse "github.com/mrvin/tasks-go/e-wallet/pkg/http/response"
)

type WalletBalance interface {
	Balance(ctx context.Context, walletID uuid.UUID) (float64, error)
}

type ResponseBalance struct {
	ID      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
	Status  string    `json:"status"`
}

func New(balanceGetter WalletBalance) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		strWalletID := req.PathValue("walletID")
		walletID, err := uuid.Parse(strWalletID)
		if err != nil {
			err := fmt.Errorf("can't get parse uuid: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		balance, err := balanceGetter.Balance(req.Context(), walletID)
		if err != nil {
			err := fmt.Errorf("get balance: %w", err)
			slog.Error(err.Error())
			if errors.Is(err, sqlstorage.ErrNoWalletID) {
				httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
				return
			}
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseBalance{
			ID:      walletID,
			Balance: balance,
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
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.Info("Wallet balance get successfully")
	}
}
