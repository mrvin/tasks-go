package send

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mrvin/tasks-go/e-wallet/internal/httpserver/handlers/wallet/send/mocks"
	"github.com/mrvin/tasks-go/e-wallet/internal/storage"
	sqlstorage "github.com/mrvin/tasks-go/e-wallet/internal/storage/sql"
)

func TestSend(t *testing.T) {
	tests := []struct {
		name         string
		walletIDFrom string
		walletIDTo   string
		amount       float64
		mockError    error
		httpStatus   int
		response     string
	}{
		{
			name:         "Success",
			walletIDFrom: "989e230a-7738-4449-8ad1-684c1f201142",
			walletIDTo:   "a2c3b089-6186-43b5-bb8b-a5b07f965168",
			amount:       12.15,
			mockError:    nil,
			httpStatus:   http.StatusOK,
			response:     `{"status":"OK"}`,
		},
		{
			name:         "No WalletID From",
			walletIDFrom: "20be130a-d4f1-458f-820e-c8c7f2fb3ce0",
			walletIDTo:   "a2c3b089-6186-43b5-bb8b-a5b07f965168",
			amount:       22.35,
			mockError:    sqlstorage.ErrNoWalletIDFrom,
			httpStatus:   http.StatusNotFound,
			response:     `{"status":"Error","error":"send transaction: no wallet-from with id"}`,
		},
		{
			name:         "No WalletID To",
			walletIDFrom: "989e230a-7738-4449-8ad1-684c1f201142",
			walletIDTo:   "8516a5d6-79db-48e0-b47a-107465e401ee",
			amount:       35.22,
			mockError:    sqlstorage.ErrNoWalletIDTo,
			httpStatus:   http.StatusBadRequest,
			response:     `{"status":"Error","error":"send transaction: no wallet-to with id"}`,
		},
	}

	mockWalletSender := mocks.NewWalletSender()

	handlerSend := New(mockWalletSender)
	for _, test := range tests {

		dataRequestSend, err := json.Marshal(RequestSend{To: uuid.Must(uuid.Parse(test.walletIDTo)), Amount: test.amount})
		if err != nil {
			t.Fatalf("cant marshal JSON: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, "/api/v1/wallet/"+test.walletIDFrom+"/send", bytes.NewReader(dataRequestSend))
		if err != nil {
			t.Fatalf("create request: %v", err)
		}

		transaction := storage.Transaction{
			WalletIDFrom: uuid.Must(uuid.Parse(test.walletIDFrom)),
			WalletIDTo:   uuid.Must(uuid.Parse(test.walletIDTo)),
			Amount:       test.amount,
		}
		mockWalletSender.On("Send", req.Context(), transaction).
			Return(test.mockError)

		res := httptest.NewRecorder()
		handlerSend(res, req)

		if res.Code != test.httpStatus {
			t.Errorf("response code is %d; want: %d", res.Code, test.httpStatus)
		}

		strResponse := res.Body.String()
		if strResponse != test.response {
			t.Errorf("response is %s; want: %s", strResponse, test.response)
		}
	}
}
