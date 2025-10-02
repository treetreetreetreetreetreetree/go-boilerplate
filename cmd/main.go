package main

import (
	"context"
	"go-boilerplate/config"
	"go-boilerplate/pkg/api"
	"go-boilerplate/pkg/database"
	"go-boilerplate/pkg/logger"
	"go-boilerplate/pkg/redis"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Setup()

	slog.Info("[APP]", "message", "Initialize app")

	cfg := config.LoadConfig(".")
	slog.Info("[APP]", "message", "current env: "+cfg.App.Env)

	var err error
	dbConnection, err := database.Setup(&cfg.Database)
	checkError(err)

	err = dbConnection.EnsureMigrations(database.Migrations)
	checkError(err)

	rdb, err := redis.Setup(context.TODO(), cfg.Redis)
	checkError(err)

	api.ServePublicServer(cfg.Server)
	api.ServeAPIDocs(cfg.Server)

	gracefulShutdown(
		func() error {
			return dbConnection.SQL.Close()
		},
		func() error {
			return rdb.Client.Close()
		},
		func() error {
			os.Exit(0)
			return nil
		},
	)
}

func gracefulShutdown(ops ...func() error) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	if <-shutdown != nil {
		for _, op := range ops {
			if err := op(); err != nil {
				slog.Error("gracefulShutdown op failed", "error", err)
				panic(err)
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
