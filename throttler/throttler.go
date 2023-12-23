package throttler

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

var ErrManyCalls = errors.New("too many calls")

type Throttler struct {
	originalTransport http.RoundTripper
	tokens            int
	muTokens          sync.Mutex
	mMethods          map[string]bool
}

func NewThrottler(
	roundTripper http.RoundTripper,
	numReq int,
	timeInterval time.Duration,
	slMethods []string,
) *Throttler {
	throttler := Throttler{
		originalTransport: roundTripper,
		tokens:            numReq,
		mMethods: map[string]bool{
			http.MethodGet:     true,
			http.MethodHead:    true,
			http.MethodPost:    true,
			http.MethodPut:     true,
			http.MethodPatch:   true,
			http.MethodDelete:  true,
			http.MethodConnect: true,
			http.MethodOptions: true,
			http.MethodTrace:   true,
		},
	}
	if slMethods != nil && len(slMethods) != 0 {
		// Делаем все методы не лиметированными
		for method := range throttler.mMethods {
			throttler.mMethods[method] = false
		}
		// Выбираем только нужные методы
		for _, method := range slMethods {
			throttler.mMethods[method] = true
		}
	}

	// Токины добавляются со скоростью numReq токинов
	// через каждый интервал времени timeInterval.
	ticker := time.NewTicker(timeInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				throttler.muTokens.Lock()
				throttler.tokens = numReq
				throttler.muTokens.Unlock()
			}
		}
	}()

	return &throttler
}

func (t *Throttler) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.mMethods[req.Method] {
		// Проверяем, имеются ли неиспользованные токены.
		t.muTokens.Lock()
		if t.tokens <= 0 {
			t.muTokens.Unlock()
			return nil, ErrManyCalls
		}
		// Уменьшаем кол-во токинов на единицу и
		// запускает RoundTrip().
		t.tokens--
		t.muTokens.Unlock()
	}
	return t.originalTransport.RoundTrip(req)
}
