package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/merch-shop/internal/logger"
	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/merch-shop/pkg/http/response"
)

type CoinSender interface {
	SendCoin(ctx context.Context, transaction *storage.Transaction) error
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"` // Имя пользователя, которому нужно отправить монеты.
	Amount uint64 `json:"amount"` // Количество монет, которые необходимо отправить.
}

func NewSendCoin(sender CoinSender) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fromUser, err := logger.GetUserNameFromCtx(req.Context())
		if err != nil {
			err := fmt.Errorf("get user name from ctx: %w", err)
			slog.ErrorContext(req.Context(), "Send coin: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		// Read json request
		var request SendCoinRequest
		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.ErrorContext(req.Context(), err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.ErrorContext(req.Context(), err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		transaction := storage.Transaction{
			FromUser: fromUser,
			ToUser:   request.ToUser,
			Amount:   request.Amount,
		}
		if err := sender.SendCoin(req.Context(), &transaction); err != nil {
			err := fmt.Errorf("send coin: %w", err)
			slog.Error(err.Error())
			if errors.Is(err, storage.ErrNoFromUser) {
				httpresponse.WriteError(res, err.Error(), http.StatusNotFound)
				return
			}
			if errors.Is(err, storage.ErrNoToUser) || errors.Is(err, storage.ErrNotEnoughFunds) {
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res, http.StatusOK)

		slog.InfoContext(req.Context(), "Transaction was successful")
	}
}
