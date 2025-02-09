package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mrvin/tasks-go/pinger/internal/storage"
)

type Conf struct {
	Addr string
}

//nolint:tagliatelle
type ResponseListHost struct {
	ListHost []storage.Host `json:"list_host"`
	Status   string         `json:"status"`
}

func ListHost(addr string) ([]storage.Host, error) {
	requestURL := "http://" + addr + "/hosts"

	resp, err := http.Get(requestURL) //nolint:gosec,bodyclose,noctx
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response status code: %d", resp.StatusCode)
	}

	var response ResponseListHost
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("unmarshal body response: %w", err)
	}

	return response.ListHost, nil
}

func CreatePing(addr string, ping *storage.Ping) error {
	requestURL := "http://" + addr + "/pings"

	requestDataJson, err := json.Marshal(ping)
	if err != nil {
		return err
	}

	resp, err := http.Post(requestURL, "application/json", bytes.NewReader(requestDataJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http response status code: %d", resp.StatusCode)
	}

	return nil
}
