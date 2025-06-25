package list

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mrvin/tasks-go/goods/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/goods/pkg/http/response"
)

const (
	defaultLimit  = 100
	defaultOffset = 0
)

type GoodLister interface {
	List(ctx context.Context, limit, offset uint64) ([]storage.Good, error)
	Meta(ctx context.Context) (int64, int64, error)
}

type ResponseMeta struct {
	Total   int64  `json:"total"`
	Removed int64  `json:"removed"`
	Limit   uint64 `json:"limit"`
	Offset  uint64 `json:"offset"`
}

type ResponseGoods struct {
	Meta   ResponseMeta   `json:"meta"`
	Goods  []storage.Good `json:"goods"`
	Status string         `json:"status"`
}

func New(lister GoodLister) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var err error
		limit := uint64(defaultLimit)
		limitStr := req.URL.Query().Get("limit")
		if limitStr != "" {
			limit, err = strconv.ParseUint(limitStr, 10, 64)
			if err != nil {
				err := fmt.Errorf("incorrect limit value: %w", err)
				slog.Error(err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
		}
		offset := uint64(defaultOffset)
		offsetStr := req.URL.Query().Get("offset")
		if offsetStr != "" {
			offset, err = strconv.ParseUint(offsetStr, 10, 64)
			if err != nil {
				err := fmt.Errorf("incorrect offset value: %w", err)
				slog.Error(err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
		}

		goods, err := lister.List(ctx, limit, offset)
		if err != nil {
			err := fmt.Errorf("get list goods: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		total, removed, err := lister.Meta(ctx)
		if err != nil {
			err := fmt.Errorf("get meta: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		meta := ResponseMeta{
			Total:   total,
			Removed: removed,
			Limit:   limit,
			Offset:  offset,
		}
		response := ResponseGoods{
			Meta:   meta,
			Goods:  goods,
			Status: "OK",
		}
		jsonResponseGoods, err := json.Marshal(response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err := res.Write(jsonResponseGoods); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.Error(err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
