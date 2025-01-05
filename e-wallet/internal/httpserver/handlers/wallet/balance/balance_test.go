package balance

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/balance/mocks"
	sqlstorage "github.com/mrvin/tasks-go/e-wallet/internal/storage/sql"
)

func TestBalance(t *testing.T) {
	tests := []struct {
		name        string
		walletID    string
		mockBalance float64
		mockError   error
		httpStatus  int
		response    string
	}{
		{
			name:        "Success",
			walletID:    "989e230a-7738-4449-8ad1-684c1f201142",
			mockBalance: 87.85,
			mockError:   nil,
			httpStatus:  http.StatusOK,
			response:    `{"id":"989e230a-7738-4449-8ad1-684c1f201142","balance":87.85,"status":"OK"}`,
		},
		{
			name:        "No wallet with id",
			walletID:    "a2c3b089-6186-43b5-bb8b-a5b07f965168",
			mockBalance: 0.0,
			mockError:   sqlstorage.ErrNoWalletID,
			httpStatus:  http.StatusNotFound,
			response:    `{"status":"Error","error":"get balance: no wallet with id"}`,
		},
	}

	mux := http.NewServeMux()
	mockBalanceGetter := mocks.NewWalletBalance()
	mux.HandleFunc(http.MethodGet+" /api/v1/wallet/"+"{walletID}", New(mockBalanceGetter))
	for _, test := range tests {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/wallet/"+test.walletID, nil)
		if err != nil {
			t.Fatalf("create request: %v", err)
		}

		mockBalanceGetter.On("Balance", req.Context(), uuid.Must(uuid.Parse(test.walletID))).
			Return(test.mockBalance, test.mockError)

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
