package handler

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewRbacEnforcer),
	fx.Provide(NewServer),
	fx.Invoke(func(server *Server, lc fx.Lifecycle) {
		server.Start(lc)
	}),
)
