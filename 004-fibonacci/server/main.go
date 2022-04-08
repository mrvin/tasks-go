package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/mrvin/tasks-go/004-fibonacci/server/config"
	pb "github.com/mrvin/tasks-go/004-fibonacci/fibpb"
	"google.golang.org/grpc"
)

var ctx = context.Background()

type serverGRPC struct {
	rdb      *redis.Client
	nMaxSave uint64
}

func main() {
	var serv serverGRPC
	var conf config.Config

	configPath := "config.yml"
	if err := conf.Parse(configPath); err != nil {
		log.Fatalf("fibserver: %v", err)
	}

	var err error
	serv.rdb, err = rdbConnect(&conf)
	if err != nil {
		log.Fatalf("fibserver: %v", err)
	}
	defer serv.rdb.Close()

	nMaxSaveStr, err := serv.rdb.Get(ctx, "nMaxSave").Result()
	if err == redis.Nil {
		log.Print("fibserver: cant get nMaxSave")
	} else {
		if err != nil {
			log.Fatalf("fibserver: %v", err)
		}
		serv.nMaxSave, err = strconv.ParseUint(nMaxSaveStr, 10, 64)
		if err != nil {
			log.Fatalf("fibserver: %v", err)
		}
	}

	if err := serv.run(&conf); err != nil {
		log.Fatalf("fibserver: %v", err)
	}
}

func (s *serverGRPC) run(conf *config.Config) error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.HostGRPC, conf.PortGRPC))
	if err != nil {
		return fmt.Errorf("сan't establish tcp connection: %v", err)
	}
	grpcServ := grpc.NewServer()
	pb.RegisterFibServer(grpcServ, s)
	if err := grpcServ.Serve(ln); err != nil {
		return fmt.Errorf("сan't run grpc server: %v", err)
	}

	return nil
}

func (s *serverGRPC) Get(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	response := make(chan *pb.Response)

	go func() {
		var slValFib []string
		nMin := req.GetNMin()
		nMax := req.GetNMax()
		if s.nMaxSave >= nMax {
			slValFib, _ = getFromCache(s.rdb, nMin, nMax)
			response <- &pb.Response{Numbers: slValFib}
		} else {
			slValFib = fib(nMax)
			response <- &pb.Response{Numbers: slValFib[nMin-1:]}
			if err := setToCache(s.rdb, slValFib, nMax); err != nil {
				log.Print(err)
				return
			}
			s.nMaxSave = nMax
		}
		close(response)
	}()

	select {
	case rsp := <-response:
		return rsp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func rdbConnect(conf *config.Config) (rdb *redis.Client, err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.DB.Host, conf.DB.Port),
		Password: conf.DB.Password,
		DB:       conf.DB.NameDB,
	})
	err = rdb.Ping(ctx).Err()

	return
}

func getFromCache(rdb *redis.Client, nMin, nMax uint64) ([]string, error) {
	slValFib := make([]string, 0, nMax-nMin)
	for i := nMin; i <= nMax; i++ {
		num, err := rdb.Get(ctx, strconv.FormatUint(i, 10)).Result()
		if err != nil {
			return nil, fmt.Errorf("can't get from cash num %d: %v", i, err)
		}
		slValFib = append(slValFib, num)
	}

	return slValFib, nil
}

func setToCache(rdb *redis.Client, slValFib []string, nMax uint64) error {
	for i := uint64(0); i < nMax; i++ {
		if err := rdb.Set(ctx, strconv.FormatUint(i+1, 10), slValFib[i], 0).Err(); err != nil {
			return fmt.Errorf("can't set to cash num %d: %v", i+1, err)
		}
	}
	if err := rdb.Set(ctx, "nMaxSave", strconv.FormatUint(nMax, 10), 0).Err(); err != nil {
		return fmt.Errorf("can't set to cash nMaxSave: %v", err)
	}

	return nil
}

func fib(nMax uint64) []string {
	a := big.NewInt(0)
	b := big.NewInt(1)

	slValFib := make([]string, 0, nMax)
	for i := uint64(0); i < nMax; i++ {
		a.Add(a, b)
		a, b = b, a
		slValFib = append(slValFib, a.String())
	}

	return slValFib
}
