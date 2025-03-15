package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/merch-shop/internal/logger"
	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/merch-shop/pkg/http/response"
)

type InfoGetter interface {
	GetBalance(ctx context.Context, userName string) (uint64, error)
	GetInventory(ctx context.Context, userName string) ([]storage.ProductQuantity, error)
	GetHistory(ctx context.Context, userName string) (storage.HistoryResponse, error)
}

type InfoResponse struct {
	Coins       uint64                    `json:"coins"`
	Inventory   []storage.ProductQuantity `json:"inventory"`
	CoinHistory storage.HistoryResponse   `json:"coinHistory"`
	Status      string                    `json:"status"`
}

func NewInfo(getter InfoGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userName, err := logger.GetUserNameFromCtx(req.Context())
		if err != nil {
			err := fmt.Errorf("get user name from ctx: %w", err)
			slog.ErrorContext(req.Context(), "Info: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		balance, err := getter.GetBalance(req.Context(), userName)
		if err != nil {
			err := fmt.Errorf("get balance: %w", err)
			slog.ErrorContext(req.Context(), "Info: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		inventory, err := getter.GetInventory(req.Context(), userName)
		if err != nil {
			err := fmt.Errorf("get inventory: %w", err)
			slog.ErrorContext(req.Context(), "Info: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		history, err := getter.GetHistory(req.Context(), userName)
		if err != nil {
			err := fmt.Errorf("get transaction history: %w", err)
			slog.ErrorContext(req.Context(), "Info: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := InfoResponse{
			Coins:       balance,
			Inventory:   inventory,
			CoinHistory: history,
			Status:      "OK",
		}
		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.ErrorContext(req.Context(), "Info: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(req.Context(), "Info: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
