package servgrpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mrvin/tasks-go/fibonacci/internal/cache"
	"github.com/mrvin/tasks-go/fibonacci/internal/fibonacci"
	"github.com/mrvin/tasks-go/fibonacci/internal/fibonacci-api"
	"google.golang.org/grpc"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ServerGRPC struct {
	cacheFib cache.Cache
}

func (s *ServerGRPC) Run(conf *Conf, cacheFib cache.Cache) error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
	if err != nil {
		return fmt.Errorf("сan't establish tcp connection: %v", err)
	}

	s.cacheFib = cacheFib
	grpcServ := grpc.NewServer()
	fibonacciapi.RegisterFibServer(grpcServ, s)
	if err := grpcServ.Serve(ln); err != nil {
		return fmt.Errorf("сan't run grpc server: %v", err)
	}

	return nil
}

func (s *ServerGRPC) Get(ctx context.Context, req *fibonacciapi.Request) (*fibonacciapi.Response, error) {
	response := make(chan *fibonacciapi.Response)

	go func() {
		from := req.GetFrom()
		to := req.GetTo()

		slValFib, err := fibonacci.GetFibNumbers(s.cacheFib, from, to)
		if err != nil {
			log.Printf("can't get fibonacci number: %v", err)
		}

		response <- &fibonacciapi.Response{Numbers: slValFib}
		close(response)
	}()

	select {
	case rsp := <-response:
		return rsp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
