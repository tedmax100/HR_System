package main

import (
	"context"
	"os"
	"os/signal"
	"syscall" // 请替换为你的实际项目路径
	"time"

	"github.com/hrdemo/internal/config"
	"github.com/hrdemo/internal/driver/db"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		config.Module,
		db.Module,
	)
	baseCtx := context.Background()
	startCtx, cancel := context.WithTimeout(baseCtx, 1*time.Minute)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.NotifyContext(baseCtx, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-c
	stopCtx, cancel := context.WithTimeout(baseCtx, 1*time.Minute)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		panic(err)
	}

}
