package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrdemo/internal/config"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ServerParams struct {
	fx.In
	Config       *config.Config
	Db           *gorm.DB
	RbacEnforcer *RbacEnforcer
	ServiceName  string `name:"service_name"`
	Version      string `name:"version"`
	Build        string `name:"build_date"`
	GoVersion    string `name:"go_version"`
}

type Server struct {
	cfg          *config.Config
	engine       *gin.Engine
	db           *gorm.DB
	srv          *http.Server
	rbacEnforcer *RbacEnforcer
	Name         string
	Version      string
	Build        string
	GoVersion    string
}

func NewServer(params ServerParams) *Server {
	r := gin.Default()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", params.Config.Port),
		Handler: r,
	}

	server := &Server{
		cfg:          params.Config,
		db:           params.Db,
		rbacEnforcer: params.RbacEnforcer,
		engine:       r,
		srv:          srv,
		Name:         params.ServiceName,
		Version:      params.Version,
		Build:        params.Build,
		GoVersion:    params.GoVersion,
	}

	server.initRoutes()
	log.Printf("server listen on %d", params.Config.Port)
	return server
}

func (s *Server) Start(lc fx.Lifecycle) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalln(err.Error())
				}
			}()

			return nil

		},
		OnStop: func(ctx context.Context) error {
			return s.srv.Shutdown(ctx)
		},
	})
	return nil
}
