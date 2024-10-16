package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/mrvin/tasks-go/url-shortener/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/url-shortener/pkg/http/response"
)

type URLCreator interface {
	CreateURL(ctx context.Context, urlToSave string, alias string) (int64, error)
}

type Request struct {
	URL   string `json:"url"             validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Alias  string `json:"alias"`
	Status string `json:"status"`
}

func NewSaveURL(creator URLCreator, defaultAliasLengthint int) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var request Request

		// Read json request
		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if request.Alias == "" {
			request.Alias = generateAlias(defaultAliasLengthint)
		}

		id, err := creator.CreateURL(req.Context(), request.URL, request.Alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				err := fmt.Errorf("alias already exists: %w", err)
				slog.InfoContext(req.Context(), "Save url: "+err.Error(), slog.String("alias", request.Alias))
				httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
				return
			}
			err := fmt.Errorf("failed save url: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error(), slog.String("alias", request.Alias))
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := Response{
			Alias:  request.Alias,
			Status: "OK",
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(req.Context(), "Save url: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.InfoContext(req.Context(), "Create new url",
			slog.Int64("id", id),
			slog.String("alias", request.Alias),
			slog.String("url", request.URL),
		)
	}
}

func generateAlias(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789")

	alias := make([]rune, length)
	for i := range alias {
		alias[i] = chars[rnd.Intn(len(chars))]
	}

	return string(alias)
}
