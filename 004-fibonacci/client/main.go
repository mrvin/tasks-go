package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/mrvin/tasks-go/004-fibonacci/fibpb"
	"google.golang.org/grpc"
)

func usage() {
	fmt.Printf("usage: %s -host hostname -port port\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var port int
	var host string
	var nMin, nMax uint64

	flag.IntVar(&port, "port", 55555, "port")
	flag.StringVar(&host, "host", "localhost", "host name")
	flag.Uint64Var(&nMin, "min", 1, "interval [min, max]")
	flag.Uint64Var(&nMax, "max", 10, "interval [min, max]")
	flag.Usage = usage
	flag.Parse()

	if nMin > nMax {
		log.Fatalf("fibclient: first arg should be less or equal second arg")
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fibclient: %v", err)
	}
	defer conn.Close()
	c := pb.NewFibClient(conn)

	req := &pb.Request{NMin: nMin, NMax: nMax}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := c.Get(ctx, req)
	if err != nil {
		log.Fatalf("fibclient: %v", err)
	}

	for i, number := range r.Numbers {
		fmt.Printf("%d - %s\n", nMin+uint64(i), number)
	}
}
