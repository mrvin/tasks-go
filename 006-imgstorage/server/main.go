//go:generate protoc -I=../api/ --go_out=../internal/imgstorageapi --go-grpc_out=require_unimplemented_servers=false:../internal/imgstorageapi ../api/imgstorage_service.proto
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrvin/tasks-go/006-imgstorage/internal/config"
	"github.com/mrvin/tasks-go/006-imgstorage/server/grpcserver"
)

func main() {
	configFile := flag.String("config", "/etc/calendar/config.yml", "path to configuration file")
	flag.Parse()

	var conf grpcserver.Config
	if err := config.Parse(*configFile, &conf); err != nil {
		log.Printf("Parse config: %v", err)
		return
	}

	err := os.MkdirAll(conf.Dir, 0750)
	if err != nil && !os.IsExist(err) {
		log.Printf("Сreate directory %s: %v", conf.Dir, err)
		return
	}

	serverGRPC, err := grpcserver.New(&conf)
	if err != nil {
		log.Printf("gRPC server: %v", err)
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT /*(Control-C)*/, syscall.SIGTERM)
	go listenForShutdown(signals, serverGRPC)

	if err := serverGRPC.Start(); err != nil {
		log.Printf("gRPC server: failed to start: %v", err)
		return
	}
}

func listenForShutdown(signals chan os.Signal, serverGRPC *grpcserver.Server) {
	<-signals
	signal.Stop(signals)

	serverGRPC.Stop()
}
