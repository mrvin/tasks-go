package main

import (
	"log"

	"github.com/mrvin/tasks-go/004-fibonacci/server/cache"
	"github.com/mrvin/tasks-go/004-fibonacci/server/config"
	"github.com/mrvin/tasks-go/004-fibonacci/server/servgrpc"
	"github.com/mrvin/tasks-go/004-fibonacci/server/servhttp"
)

func main() {
	var conf config.Config
	configPath := "config.yml"
	if err := conf.Parse(configPath); err != nil {
		log.Fatalf("fibserver: %v", err)
	}

	var cacheFib cache.CacheRDB
	if err := cacheFib.Connect(&conf); err != nil {
		log.Fatalf("fibserver: %v", err)
	}

	done := make(chan struct{})
	go func() {
		var serv servgrpc.ServerGRPC
		if err := serv.Run(&conf, &cacheFib); err != nil {
			log.Printf("fibserver: %v", err)
		}
		done <- struct{}{}
	}()

	var serv servhttp.ServerHTTP
	if err := serv.Run(&conf, &cacheFib); err != nil {
		log.Printf("fibserver: %v", err)
	}

	<-done
	cacheFib.Close()
}
