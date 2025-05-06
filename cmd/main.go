package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/michelsazevedo/hawkeye/config"
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
		fx.Invoke(config.Routes),
	)

	app.Run()
}
