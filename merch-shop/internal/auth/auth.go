package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/mrvin/tasks-go/merch-shop/internal/logger"
)

//nolint:tagliatelle
type Conf struct {
	SecretKey           string        `yaml:"secret_key"`
	TokenValidityPeriod time.Duration `yaml:"token_validity_period"`
}

type AuthService struct {
	SecretKey           string
	TokenValidityPeriod time.Duration
}

func New(conf *Conf) *AuthService {
	return &AuthService{conf.SecretKey, conf.TokenValidityPeriod}
}

func (a *AuthService) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"iat":      time.Now().Unix(),                            // IssuedAt
			"exp":      time.Now().Add(a.TokenValidityPeriod).Unix(), // ExpiresAt
		},
	)
	tokenString, err := token.SignedString([]byte(a.SecretKey))
	if err != nil {
		return "", fmt.Errorf("create token: %w", err)
	}

	return tokenString, nil
}

func (a *AuthService) Auth(next http.HandlerFunc) http.HandlerFunc {
	handler := func(res http.ResponseWriter, req *http.Request) {
		authHeaderValue := req.Header.Get("Authorization")
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeaderValue, bearerPrefix) {
			http.Error(res, "request does not contain an authorization bearer token", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeaderValue, bearerPrefix)

		claims, err := a.parseToken(tokenString)
		if err != nil {
			http.Error(res, err.Error(), http.StatusUnauthorized)
			return
		}
		userName := claims["username"]

		ctx := logger.WithUserName(req.Context(), userName.(string))

		next(res, req.WithContext(ctx)) // Pass request to next handler
	}

	return http.HandlerFunc(handler)
}

func (a *AuthService) parseToken(tokenStr string) (jwt.MapClaims, error) {
	// Validate token.
	token, err := jwt.Parse(
		tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.SecretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
