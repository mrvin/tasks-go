package spelling

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

const requestTimeout = 10

func requestBuildCheckText(text string) (string, error) {
	const requestURL = "https://speller.yandex.net/services/spellservice.json/checkText?text=синхрафазатрон+в+дубне"

	url, err := url.Parse(requestURL)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}

	values := url.Query()
	values.Set("text", text)
	url.RawQuery = values.Encode()

	return url.String(), nil
}

func Check(ctx context.Context, text string) (bool, error) {
	requestURL, err := requestBuildCheckText(text)
	if err != nil {
		return false, fmt.Errorf("check text: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, requestTimeout*time.Second)
	defer cancel()

	// Create a new request using http
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return false, fmt.Errorf("create request: %w", err)
	}

	slog.Info("Check text", slog.String("url", requestURL))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("http response status code: %d", resp.StatusCode)
	}

	var respData []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return false, fmt.Errorf("unmarshal body response: %w", err)
	}

	if len(respData) != 0 {
		fmt.Println(respData)
		return false, nil
	}

	return true, nil
}
