package servgrpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mrvin/tasks-go/004-fibonacci/fibpb"
	"github.com/mrvin/tasks-go/004-fibonacci/server/cache"
	"github.com/mrvin/tasks-go/004-fibonacci/server/config"
	"github.com/mrvin/tasks-go/004-fibonacci/server/fibonacci"
	"google.golang.org/grpc"
)

type ServerGRPC struct {
	cacheFib cache.Cache
}

func (s *ServerGRPC) Run(conf *config.Config, cacheFib cache.Cache) error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.HostGRPC, conf.PortGRPC))
	if err != nil {
		return fmt.Errorf("сan't establish tcp connection: %v", err)
	}

	s.cacheFib = cacheFib
	grpcServ := grpc.NewServer()
	fibpb.RegisterFibServer(grpcServ, s)
	if err := grpcServ.Serve(ln); err != nil {
		return fmt.Errorf("сan't run grpc server: %v", err)
	}

	return nil
}

func (s *ServerGRPC) Get(ctx context.Context, req *fibpb.Request) (*fibpb.Response, error) {
	response := make(chan *fibpb.Response)

	go func() {
		from := req.GetFrom()
		to := req.GetTo()

		slValFib, err := fibonacci.GetFibNumbers(s.cacheFib, from, to)
		if err != nil {
			log.Printf("can't get fibonacci number: %v", err)
		}

		response <- &fibpb.Response{Numbers: slValFib}
		close(response)
	}()

	select {
	case rsp := <-response:
		return rsp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
