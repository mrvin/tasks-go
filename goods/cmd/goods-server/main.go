package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/mrvin/tasks-go/goods/internal/config"
	"github.com/mrvin/tasks-go/goods/internal/httpserver"
	"github.com/mrvin/tasks-go/goods/internal/logger"
	natsmq "github.com/mrvin/tasks-go/goods/internal/queue/nats"
	"github.com/mrvin/tasks-go/goods/internal/storage"
	chstorage "github.com/mrvin/tasks-go/goods/internal/storage/clickhouse"
	sqlstorage "github.com/mrvin/tasks-go/goods/internal/storage/sql"
)

const BatchSize = 2

func main() {
	// init config
	var conf config.Config
	conf.LoadFromEnv()

	logFile, err := logger.Init(&conf.Logger)
	if err != nil {
		log.Printf("Init logger: %v\n", err)
		return
	}
	slog.Info("Init logger", slog.String("Logging level", conf.Logger.Level))
	defer func() {
		if err := logFile.Close(); err != nil {
			slog.Error("Close log file: " + err.Error())
		}
	}()

	ctx := context.Background()
	slog.Info("Storage in sql database")
	sqlStorage, err := sqlstorage.New(ctx, &conf.Postgres)
	if err != nil {
		slog.Error("Failed to init sql storage: " + err.Error())
		return
	}
	slog.Info("Connected to sql database")
	defer func() {
		if err := sqlStorage.Close(); err != nil {
			slog.Error("Failed to close sql storage: " + err.Error())
		} else {
			slog.Info("Closing the sql database connection")
		}
	}()

	// Connected to ClickHouse
	slog.Info("Storage in clickhouse")
	chStorage, err := chstorage.New(ctx, &conf.ClickHouse)
	if err != nil {
		slog.Error("Failed to init clickhouse storage: " + err.Error())
		return
	}
	slog.Info("Connected to clickhouse database")
	defer func() {
		if err := chStorage.Close(); err != nil {
			slog.Error("Failed to close clickhouse storage: " + err.Error())
		} else {
			slog.Info("Closing the clickhouse database connection")
		}
	}()

	mq, err := natsmq.New(&conf.MQ)
	if err != nil {
		slog.Error("Failed to init mq: " + err.Error())
		return
	}
	slog.Info("Connected to NATS")
	defer func() {
		if err := mq.Close(); err != nil {
			slog.Error("Failed to close NATS: " + err.Error())
		} else {
			slog.Info("Closing the NATS connection")
		}
	}()

	server := httpserver.New(&conf.HTTP, sqlStorage, mq)

	go func() {
		batch := make([]storage.Event, 0, BatchSize)
		for event := range mq.EventsCh {
			batch = append(batch, event)
			if len(batch) >= BatchSize {
				if err := chStorage.BatchInsert(ctx, batch); err != nil {
					slog.Error("Batch insert: " + err.Error())
				}
				batch = batch[:0]
			}
		}
		if len(batch) != 0 {
			if err := chStorage.BatchInsert(ctx, batch); err != nil {
				slog.Error("Batch insert: " + err.Error())
			}
		}
		slog.Info("Stop inserting events")
	}()

	server.Run(ctx)
}
