package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/merch-shop/internal/logger"
	httpresponse "github.com/mrvin/tasks-go/merch-shop/pkg/http/response"
)

type ProductBuyer interface {
	BuyProduct(ctx context.Context, userName, productName string) error
}

func NewBuyProduct(buyer ProductBuyer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userName, err := logger.GetUserNameFromCtx(req.Context())
		if err != nil {
			err := fmt.Errorf("get user name from ctx: %w", err)
			slog.ErrorContext(req.Context(), "Buy product: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		productName := req.PathValue("productName")
		if productName == "" {
			err := errors.New("empty product name")
			slog.ErrorContext(req.Context(), "Buy product: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := buyer.BuyProduct(req.Context(), userName, productName); err != nil {
			slog.ErrorContext(req.Context(), "Buy product: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		httpresponse.WriteOK(res, http.StatusOK)

		slog.InfoContext(req.Context(), "The purchase was successful")
	}
}
