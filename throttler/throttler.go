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
}

func NewThrottler(roundTripper http.RoundTripper, numReq int, timeInterval time.Duration) *Throttler {
	throttler := Throttler{
		originalTransport: roundTripper,
		tokens:            numReq,
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

	return t.originalTransport.RoundTrip(req)
}
