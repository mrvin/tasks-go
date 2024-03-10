package getidenticon

import (
	"fmt"
	"io/ioutil"
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
		if err != redis.Nil {
			if err != nil {
				slog.Error(err.Error())
				return
			}
			fmt.Println("Get from cache")
			image = []byte(imageStr)
		} else {
			fmt.Println("Cache miss")
			url := queryBuildDnmonster(name, 250, dnmonsterAddr)

			fmt.Printf("URL: %s\n", url)

			resp, err := http.Get(url)
			if err != nil {
				slog.Error("Get: " + err.Error())
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				slog.Error("Response status", slog.Int("code", resp.StatusCode), slog.String("url", url))
				return
			}
			image, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				slog.Error("ReadAll: " + err.Error())
				return
			}

			if err := cache.Set(req.Context(), name, image, 0).Err(); err != nil {
				slog.Error("can't set to cache val " + name + ":" + err.Error())
				return
			}
		}

		res.Write(image)
		res.Header().Set("Content-Type", "image/png")
	}
}

func queryBuildDnmonster(name string, size int, confDnmonsterStr string) string {
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
