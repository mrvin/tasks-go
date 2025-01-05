package create

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/create/mocks"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name         string
		balance      float64
		mockWalletID string
		mockError    error
		httpStatus   int
		response     string
	}{
		{
			name:         "Success",
			balance:      100,
			mockWalletID: "989e230a-7738-4449-8ad1-684c1f201142",
			mockError:    nil,
			httpStatus:   http.StatusCreated,
			response:     `{"id":"989e230a-7738-4449-8ad1-684c1f201142","balance":100,"status":"OK"}`,
		},
	}

	mux := http.NewServeMux()
	mockWalletCreator := mocks.NewWalletCreator()
	mux.HandleFunc(http.MethodPost+" /api/v1/wallet", New(mockWalletCreator))
	for _, test := range tests {
		req, err := http.NewRequest(http.MethodPost, "/api/v1/wallet", nil)
		if err != nil {
			t.Fatalf("create request: %v", err)
		}

		mockWalletCreator.On("Create", req.Context(), test.balance).
			Return(uuid.Must(uuid.Parse(test.mockWalletID)), test.mockError)

		res := httptest.NewRecorder()
		mux.ServeHTTP(res, req)

		if res.Code != test.httpStatus {
			t.Errorf("response code is %d; want: %d", res.Code, test.httpStatus)
		}

		strResponse := res.Body.String()
		if strResponse != test.response {
			t.Errorf("response is %s; want: %s", strResponse, test.response)
		}
	}
}
