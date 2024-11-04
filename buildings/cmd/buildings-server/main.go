package main

import (
	"context"
	"flag"
	"log"

	"github.com/mrvin/tasks-go/buildings/internal/config"
	"github.com/mrvin/tasks-go/buildings/internal/httpserver"
	sqlstorage "github.com/mrvin/tasks-go/buildings/internal/storage/sql"
)

type Config struct {
	DB   sqlstorage.Conf `yaml:"db"`
	HTTP httpserver.Conf `yaml:"http"`
}

func main() {
	log.Println("Starting buildings server")
	configFile := flag.String("config", "/etc/buildings/buildings.yml", "path to configuration file")
	flag.Parse()

	// Init config
	log.Println("Initializing configuration")
	var conf Config
	if err := config.Parse(*configFile, &conf); err != nil {
		log.Printf("Failed parse config: %v", err)
		return
	}

	// Init storage
	log.Println("Initializing storage")
	ctx := context.Background()
	st, err := sqlstorage.New(ctx, &conf.DB)
	if err != nil {
		log.Printf("Failed to init storage: %v", err)
		return
	}
	defer st.Close()
	log.Printf("Connected to database")

	// Init HTTP server
	log.Println("Initializing HTTP server")
	server := httpserver.New(st)

	// Start HTTP server
	server.Start(&conf.HTTP)
}
