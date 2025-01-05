package withdraw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	httpresponse "github.com/mrvin/tasks-go/e-wallet/pkg/http/response"
)

type WalletWithdrawer interface {
	Withdraw(ctx context.Context, walletID uuid.UUID, amount float64) error
}

type RequestWithdraw struct {
	Amount float64 `json:"amount"`
}

func New(withdrawer WalletWithdrawer, minimalAmount float64) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		strWalletID := req.PathValue("walletID")
		walletIDFrom, err := uuid.Parse(strWalletID)
		if err != nil {
			err := fmt.Errorf("can't get parse uuid: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		// Read json request
		var request RequestWithdraw

		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if request.Amount < minimalAmount {
			err := errors.New("amount is too small")
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := withdrawer.Withdraw(req.Context(), walletIDFrom, request.Amount); err != nil {
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res, http.StatusOK)

		slog.Info("Funds were debited successfully")
	}
}
