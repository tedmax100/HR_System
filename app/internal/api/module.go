package api

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewRbacEnforcer),
	fx.Provide(NewServer),
)
