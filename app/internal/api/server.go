package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrdemo/internal/config"
	"go.uber.org/fx"
)

type ServerParams struct {
	fx.In
	Config       *config.Config
	RbacEnforcer *RbacEnforcer
}

type Server struct {
	cfg          *config.Config
	engine       *gin.Engine
	srv          *http.Server
	rbacEnforcer *RbacEnforcer
}

func NewServer(params ServerParams) *Server {
	r := gin.Default()

	return &Server{
		cfg:          params.Config,
		rbacEnforcer: params.RbacEnforcer,
		engine:       r,
	}
}

func (s Server) Start(lc fx.Lifecycle) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return err
			}
			return nil

		},
		OnStop: func(ctx context.Context) error {
			return s.srv.Shutdown(ctx)
		},
	})
	return nil
}
