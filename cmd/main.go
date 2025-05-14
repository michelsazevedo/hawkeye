package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/michelsazevedo/hawkeye/internal/config"
	"github.com/michelsazevedo/hawkeye/internal/modules"
	"github.com/rs/zerolog/log"
)

func main() {
	app := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return config.NewFxZerolog(log.Logger)
		}),
		fx.Provide(
			config.NewConfig,
			config.NewApplication,
		),
		modules.Modules(),
		fx.Invoke(config.Routes),
	)

	app.Run()
}
