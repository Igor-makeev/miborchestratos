package main

import (
	"context"
	"fmt"
	config "miborchestrator/configs"
	"miborchestrator/internal/handler"
	"miborchestrator/internal/repository"
	"miborchestrator/internal/service"
	"miborchestrator/pkg/http/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}
	ctx := context.Background()
	cfg := config.NewConfig()
	pc, err := repository.NewPostgresClient(cfg)
	if err != nil {
		logrus.Fatalf("error:%v")
	}
	repo := repository.NewRepository(pc)
	service := service.NewService(ctx, cfg, repo)
	service.WcQueue.Run(ctx)
	handler := handler.NewHandler(service)

	srv := new(server.Server)

	serverErrChan := srv.Run(cfg.ServerPort, handler)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-signals:

		fmt.Println("main: got terminate signal. Shutting down...")

		if err := srv.Shutdown(); err != nil {
			fmt.Printf("main: received an error while shutting down the server: %v", err)
		}

		if err := service.Close(ctx); err != nil {
			fmt.Printf("main: an error was received while closing the service: %v", err)
		}

	case <-serverErrChan:
		fmt.Println("main: got server err signal. Shutting down...")

		if err := service.Close(ctx); err != nil {
			fmt.Printf("main: an error was received while closing the service: %v", err)
		}
	}
}
