package history

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
	sqlstorage "github.com/mrvin/tasks-go/e-wallet/internal/storage/sql"
	httpresponse "github.com/mrvin/tasks-go/e-wallet/pkg/http/response"
)

type WalletHistory interface {
	HistoryTransactions(ctx context.Context, wallet uuid.UUID) ([]storage.Transaction, error)
}

type ResponseHistory struct {
	Transactions []storage.Transaction `json:"transactions"`
	Status       string                `json:"status"`
}

func New(historyGetter WalletHistory) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		strWalletID := req.URL.Path[15:51]
		walletID, err := uuid.Parse(strWalletID)
		if err != nil {
			err := fmt.Errorf("can't get parse uuid: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		historyTransactions, err := historyGetter.HistoryTransactions(req.Context(), walletID)
		if err != nil {
			err := fmt.Errorf("get history transactions: %w", err)
			slog.Error(err.Error())
			if errors.Is(err, sqlstorage.ErrNoWalletID) {
				httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
				return
			}
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := ResponseHistory{
			Transactions: historyTransactions,
			Status:       "OK",
		}

		jsonResponseTransactions, err := json.Marshal(response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponseTransactions); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
