package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hrdemo/internal/config"
	"github.com/hrdemo/internal/driver/db"
	"github.com/hrdemo/internal/handler"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

var (
	Name       = "HR_System"
	GoVersiuon = "111"
	GitCommit  = "aaa"
	BuildDate  = "1235"
)

func main() {
	app := fx.New(
		fx.Provide(fx.Annotated{
			Name:   "service_name",
			Target: func() string { return Name },
		}),
		fx.Provide(fx.Annotated{
			Name:   "version",
			Target: func() string { return GitCommit },
		}),
		fx.Provide(fx.Annotated{
			Name:   "build_date",
			Target: func() string { return BuildDate },
		}),
		fx.Provide(fx.Annotated{
			Name:   "go_version",
			Target: func() string { return GoVersiuon },
		}),
		config.Module,
		db.Module,
		handler.Module,
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
