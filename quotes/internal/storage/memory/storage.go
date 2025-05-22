package memory

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/mrvin/tasks-go/quotes/internal/storage"
)

type Storage struct {
	maxQuoteID int64
	mQuotes    map[int64]storage.QuoteWithoutID
	sync.RWMutex
}

func New() *Storage {
	var s Storage
	s.mQuotes = make(map[int64]storage.QuoteWithoutID)

	return &s
}

func (s *Storage) Create(_ context.Context, quote *storage.QuoteWithoutID) (int64, error) {
	s.Lock()
	defer s.Unlock()

	s.maxQuoteID++
	s.mQuotes[s.maxQuoteID] = *quote

	return s.maxQuoteID, nil
}

func (s *Storage) GetRandom(_ context.Context) (*storage.Quote, error) {
	if len(s.mQuotes) == 0 {
		return nil, storage.ErrEmptyStorage
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	rndNum := rnd.Intn(len(s.mQuotes))

	s.RLock()
	defer s.RUnlock()

	count := 0
	var resultQuote storage.Quote
	for id, quote := range s.mQuotes {
		if count == rndNum {
			resultQuote.ID = id
			resultQuote.Author = quote.Author
			resultQuote.Text = quote.Text
			break
		}
		count++
	}

	return &resultQuote, nil
}

func (s *Storage) Delete(_ context.Context, id int64) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.mQuotes[id]; !ok {
		return fmt.Errorf("%w: %d", storage.ErrNoQuoteID, id)
	}
	delete(s.mQuotes, id)

	return nil
}

func (s *Storage) List(_ context.Context, author string) ([]storage.Quote, error) {
	s.RLock()
	defer s.RUnlock()

	resultList := make([]storage.Quote, 0, len(s.mQuotes))
	for id, quote := range s.mQuotes {
		if author == "" || author == quote.Author {
			resultList = append(resultList,
				storage.Quote{
					ID:     id,
					Author: quote.Author,
					Text:   quote.Text,
				})
		}
	}

	return resultList, nil
}
