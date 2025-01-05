package history

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/history/mocks"
	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
	sqlstorage "github.com/mrvin/tasks-go/e-wallet/internal/storage/sql"
)

func TestHistory(t *testing.T) {
	tests := []struct {
		name             string
		walletID         string
		mockTransactions []storage.Transaction
		mockError        error
		httpStatus       int
		response         string
	}{
		{
			name:     "Success",
			walletID: "989e230a-7738-4449-8ad1-684c1f201142",
			mockTransactions: []storage.Transaction{
				{
					Time:         time.Date(2022, time.February, 27, 10, 0, 0, 0, time.UTC),
					WalletIDFrom: uuid.Must(uuid.Parse("989e230a-7738-4449-8ad1-684c1f201142")),
					WalletIDTo:   uuid.Must(uuid.Parse("a2c3b089-6186-43b5-bb8b-a5b07f965168")),
					Amount:       12.15,
				},
			},
			mockError:  nil,
			httpStatus: http.StatusOK,
			response:   `{"transactions":[{"time":"2022-02-27T10:00:00Z","from":"989e230a-7738-4449-8ad1-684c1f201142","to":"a2c3b089-6186-43b5-bb8b-a5b07f965168","amount":12.15}],"status":"OK"}`,
		},
		{
			name:             "No WalletID",
			walletID:         "ecdd11f2-a871-4d52-be41-ff7a3fbb7f5b",
			mockTransactions: nil,
			mockError:        sqlstorage.ErrNoWalletID,
			httpStatus:       http.StatusNotFound,
			response:         `{"status":"Error","error":"get history transactions: no wallet with id"}`,
		},
	}

	mux := http.NewServeMux()
	mockWalletHistory := mocks.NewWalletHistory()
	mux.HandleFunc(http.MethodGet+" /api/v1/wallet/{walletID}/history", New(mockWalletHistory))
	for _, test := range tests {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/wallet/"+test.walletID+"/history", nil)
		if err != nil {
			t.Fatalf("create request: %v", err)
		}
		mockWalletHistory.On("HistoryTransactions", req.Context(), uuid.Must(uuid.Parse(test.walletID))).
			Return(test.mockTransactions, test.mockError)

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
