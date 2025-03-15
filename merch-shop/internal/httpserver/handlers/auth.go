package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mrvin/tasks-go/merch-shop/internal/app"
	"github.com/mrvin/tasks-go/merch-shop/internal/auth"
	"github.com/mrvin/tasks-go/merch-shop/internal/storage"
	httpresponse "github.com/mrvin/tasks-go/merch-shop/pkg/http/response"
	"golang.org/x/crypto/bcrypt"
)

type Account interface {
	GetAccount(ctx context.Context, userName string) (*storage.Account, error)
	CreateAccount(ctx context.Context, userName, hashPassword string, startingBalance uint64) error
}

type AuthRequest struct {
	Username string `json:"username"` // Имя пользователя для аутентификации.
	Password string `json:"password"` // Пароль для аутентификации.
}

type AuthResponse struct {
	Token  string `json:"token"` // JWT-токен для доступа к защищенным ресурсам.
	Status string `json:"status"`
}

func NewAuth(conf *app.Conf, account Account, a *auth.AuthService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Read json request
		var request AuthRequest
		body, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			err := fmt.Errorf("read body request: %w", err)
			slog.ErrorContext(req.Context(), "Auth: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &request); err != nil {
			err := fmt.Errorf("unmarshal body request: %w", err)
			slog.ErrorContext(req.Context(), "Auth: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		acc, err := account.GetAccount(req.Context(), request.Username)
		if err != nil {
			if errors.Is(err, storage.ErrAccountNotFound) {
				hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
				if err != nil {
					err := fmt.Errorf("generate hash password: %w", err)
					slog.ErrorContext(req.Context(), "Auth: "+err.Error())
					httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
					return
				}
				if err := account.CreateAccount(req.Context(), request.Username, string(hashPassword), conf.StartingBalance); err != nil {
					err := fmt.Errorf("create new account: %w", err)
					slog.ErrorContext(req.Context(), "Auth: "+err.Error())
					httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
					return
				}
				slog.InfoContext(req.Context(), "Create new account",
					slog.String("username", request.Username),
				)
			} else {
				err := fmt.Errorf("get account: %w", err)
				slog.ErrorContext(req.Context(), "Auth: "+err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if acc != nil {
			if err := bcrypt.CompareHashAndPassword([]byte(acc.HashPassword), []byte(request.Password)); err != nil {
				err := fmt.Errorf("compare hash and password: %w", err)
				slog.ErrorContext(req.Context(), "Auth: "+err.Error())
				httpresponse.WriteError(res, err.Error(), http.StatusUnauthorized)
				return
			}
		}
		tokenString, err := a.CreateToken(request.Username)
		if err != nil {
			slog.ErrorContext(req.Context(), "Auth: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write json response
		response := AuthResponse{
			Token:  tokenString,
			Status: "OK",
		}
		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			err := fmt.Errorf("marshal response: %w", err)
			slog.ErrorContext(req.Context(), "Auth: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		if acc == nil {
			res.WriteHeader(http.StatusCreated)
		} else {
			res.WriteHeader(http.StatusOK)
		}
		if _, err := res.Write(jsonResponse); err != nil {
			err := fmt.Errorf("write response: %w", err)
			slog.ErrorContext(req.Context(), "Auth: "+err.Error())
			httpresponse.WriteError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		slog.InfoContext(req.Context(), "Create token",
			slog.String("username", request.Username),
		)
	}
}
