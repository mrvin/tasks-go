package send

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
	sqlstorage "github.com/mrvin/tasks-go/e-wallet/internal/storage/sql"
	httpresponse "github.com/mrvin/tasks-go/e-wallet/pkg/http/response"
)

type WalletSender interface {
	Send(ctx context.Context, transaction storage.Transaction) error
}

type RequestSend struct {
	To     uuid.UUID `json:"to"`
	Amount float64   `json:"amount"`
}

func New(sender WalletSender) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		strWalletID := req.URL.Path[15:51]
		walletIDFrom, err := uuid.Parse(strWalletID)
		if err != nil {
			err := fmt.Errorf("can't get parse uuid: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		// Read json request
		var request RequestSend

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

		if request.Amount < 0.01 {
			err := errors.New("amount is too small")
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		transaction := storage.Transaction{
			WalletIDFrom: walletIDFrom,
			WalletIDTo:   request.To,
			Amount:       request.Amount,
		}

		if err := sender.Send(req.Context(), transaction); err != nil {
			err := fmt.Errorf("send transaction: %w", err)
			slog.Error(err.Error())
			if errors.Is(err, sqlstorage.ErrNoWalletIDFrom) {
				httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
				return
			}
			if errors.Is(err, sqlstorage.ErrNoWalletIDTo) || errors.Is(err, sqlstorage.ErrNotEnoughFunds) {
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res)

		slog.Info("Transaction was successful")
	}
}
