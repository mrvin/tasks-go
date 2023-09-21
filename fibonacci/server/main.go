package main

import (
	"log"

	"github.com/mrvin/tasks-go/004-fibonacci/server/cache"
	"github.com/mrvin/tasks-go/004-fibonacci/server/config"
	"github.com/mrvin/tasks-go/004-fibonacci/server/servgrpc"
	"github.com/mrvin/tasks-go/004-fibonacci/server/servhttp"
)

func main() {
	defer trace()()

	log.Println("Configuration...")
	var conf config.Config
	configPath := "config.yml"
	if err := conf.Parse(configPath); err != nil {
		log.Fatalf("fibserver: %v", err)
	}

	log.Println("Connecting to db and initializing cache...")
	var cacheFib cache.CacheRDB
	if err := cacheFib.Connect(&conf); err != nil {
		log.Fatalf("fibserver: %v", err)
	}

	done := make(chan struct{})
	go func() {
		log.Println("Start grpc server")
		var serv servgrpc.ServerGRPC
		if err := serv.Run(&conf, &cacheFib); err != nil {
			log.Printf("fibserver: %v", err)
		}
		done <- struct{}{}
	}()

	log.Println("Start http server")
	var serv servhttp.ServerHTTP
	if err := serv.Run(&conf, &cacheFib); err != nil {
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
