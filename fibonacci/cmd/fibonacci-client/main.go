package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mrvin/tasks-go/fibonacci/internal/fibonacci-api"
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
	var from, to uint64

	flag.IntVar(&port, "port", 55555, "port")
	flag.StringVar(&host, "host", "localhost", "host name")
	flag.Uint64Var(&from, "from", 1, "interval [min, max]")
	flag.Uint64Var(&to, "to", 10, "interval [min, max]")
	flag.Usage = usage
	flag.Parse()

	if from > to {
		log.Fatalf("fibclient: first arg should be less or equal second arg")
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fibclient: %v", err)
	}
	defer conn.Close()
	c := fibonacciapi.NewFibClient(conn)

	req := &fibonacciapi.Request{From: from, To: to}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := c.Get(ctx, req)
	if err != nil {
		log.Fatalf("fibclient: %v", err)
	}

	for i, number := range r.Numbers {
		fmt.Printf("%d - %s\n", from+uint64(i), number)
	}
}
