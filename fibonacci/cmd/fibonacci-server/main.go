package main

import (
	"flag"
	"log"

	"github.com/mrvin/tasks-go/fibonacci/internal/cache"
	"github.com/mrvin/tasks-go/fibonacci/internal/config"
	"github.com/mrvin/tasks-go/fibonacci/internal/servgrpc"
	"github.com/mrvin/tasks-go/fibonacci/internal/servhttp"
)

type Config struct {
	GRPC servgrpc.Conf `yaml:"grpc"`
	HTTP servhttp.Conf `yaml:"http"`
	DB   cache.Conf    `yaml:"db"`
}

func main() {
	defer trace()()

	log.Println("Configuration...")
	configFile := flag.String("config", "/etc/calendar/fibonacci-server.yml", "path to configuration file")
	flag.Parse()
	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		log.Fatalf("fibserver: %v", err)
	}

	log.Println("Connecting to db and initializing cache...")
	var cacheFib cache.CacheRDB
	if err := cacheFib.Connect(&conf.DB); err != nil {
		log.Fatalf("fibserver: %v", err)
	}

	done := make(chan struct{})
	go func() {
		log.Println("Start grpc server")
		var serv servgrpc.ServerGRPC
		if err := serv.Run(&conf.GRPC, &cacheFib); err != nil {
			log.Printf("fibserver: %v", err)
		}
		done <- struct{}{}
	}()

	log.Println("Start http server")
	var serv servhttp.ServerHTTP
	if err := serv.Run(&conf.HTTP, &cacheFib); err != nil {
		log.Printf("fibserver: %v", err)
	}

	<-done
	cacheFib.Close()
}

func trace() func() {
	log.Println("Start fibonacci server")

	return func() {
		log.Println("Stop fibonacci server")
	}
}
