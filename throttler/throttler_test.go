package throttler

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

const googleURL = "https://www.google.com/"

func TestThrottler(t *testing.T) {
	throttled := NewThrottler(http.DefaultTransport, 2, 10*time.Second)

	client := http.Client{
		Transport: throttled,
	}

	resp, err := client.Get(googleURL)
	if err != nil {
		t.Fatalf("%v", err)
	}
	resp.Body.Close()
	resp, err = client.Get(googleURL)
	if err != nil {
		t.Fatalf("%v", err)
	}
	resp.Body.Close()
	resp, err = client.Get(googleURL)
	if !errors.Is(err, ErrManyCalls) {
		t.Fatalf("%v", err)
	}
	time.Sleep(10 * time.Second)
	resp, err = client.Get(googleURL)
	if err != nil {
		t.Fatalf("%v", err)
	}
	resp.Body.Close()
}
