package servhttp

import (
//	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/mrvin/tasks-go/004-fibonacci/server/cache"
	"github.com/mrvin/tasks-go/004-fibonacci/server/config"
	"github.com/mrvin/tasks-go/004-fibonacci/server/fibonacci"
)

type ServerHTTP struct {
	pr       *pathResolver
	cacheFib cache.Cache
}

type Numbers struct {
	Numbers []string `json:"numbers"`
}

func (s *ServerHTTP) Run(conf *config.Config, cacheFib cache.Cache) error {
	s.pr = newPathResolver()
	s.pr.Add("GET /fibonacci", fibonaccihHandler)
	s.cacheFib = cacheFib
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.HostHTTP, conf.PortHTTP), s); err != nil {
		return err
	}

	return nil
}

func (s *ServerHTTP) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path

	for pattern, handlerFunc := range s.pr.handlers {
		if ok, err := path.Match(pattern, check); ok && err == nil {
			handlerFunc(res, req, s.cacheFib)
			return
		} else if err != nil {
			fmt.Fprint(res, err)
		}
	}

	http.NotFound(res, req)
}

type pathResolver struct {
	handlers map[string]func(res http.ResponseWriter, req *http.Request, cacheFib cache.Cache)
}

func newPathResolver() *pathResolver {
	return &pathResolver{make(map[string]func(res http.ResponseWriter, req *http.Request, cacheFib cache.Cache))}
}

func (p *pathResolver) Add(path string, handler func(res http.ResponseWriter, req *http.Request, cacheFib cache.Cache)) {
	p.handlers[path] = handler
}

func fibonaccihHandler(res http.ResponseWriter, req *http.Request, cacheFib cache.Cache) {
	query := req.URL.Query()
	fromStr := query.Get("from")
	from, err := strconv.ParseUint(fromStr, 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	toStr := query.Get("to")
	to, err := strconv.ParseUint(toStr, 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if from > to {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	slValFib, err := fibonacci.GetFibNumbers(cacheFib, from, to)
	if err != nil {
		log.Printf("can't get fibonacci number: %v", err)
		return
	}

//	resp := Numbers{slValFib}
//	jsonResp, err := json.Marshal(resp)
//	if err != nil {
//		log.Printf("can't marshaling json: %v", err)
//	}
//	res.Header().Set("Content-Type", "application/json")

	fmt.Fprint(res, slValFib)
}
