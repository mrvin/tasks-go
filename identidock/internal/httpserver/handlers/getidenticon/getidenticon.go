package getidenticon

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func New(dnmonsterAddr string, cache *redis.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var image []byte

		name := req.URL.Query().Get("name")

		imageStr, err := cache.Get(req.Context(), name).Result()
		if !errors.Is(err, redis.Nil) {
			if err != nil {
				slog.Error(err.Error())
				return
			}
			slog.Info("Get from cache")
			image = []byte(imageStr)
		} else {
			url := queryBuildDnmonster(name, 250, dnmonsterAddr)
			slog.Info("Cache miss", slog.String("URL", url))

			image, err = getImage(req.Context(), url)
			if err != nil {
				slog.Error("Can't get image: " + err.Error())
				return
			}
			if err := cache.Set(req.Context(), name, image, 0).Err(); err != nil {
				slog.Error("Can't set to cache val " + name + ":" + err.Error())
				return
			}
		}

		imageSize, err := res.Write(image)
		if err != nil {
			slog.Error("Can't write image: " + err.Error())
			return
		}
		res.Header().Set("Content-Type", "image/png")

		slog.Info("Image downloaded", slog.Int("bytes", imageSize))
	}
}

func queryBuildDnmonster(name string, size int, confDnmonsterStr string) string {
	//nolint:exhaustruct
	query := url.URL{
		Scheme: "http",
		Host:   confDnmonsterStr,
		Path:   "monster/" + name,
	}

	val := url.Values{
		"size": {"250"},
	}

	if size > 0 {
		val.Set("size", strconv.Itoa(size))
	}

	query.RawQuery = val.Encode()

	return query.String()
}

func getImage(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http new request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http response status code: %d", resp.StatusCode)
	}
	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all http response body: %w", err)
	}

	return image, nil
}
