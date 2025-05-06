package config

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"

	m "github.com/michelsazevedo/hawkeye/middleware"
)

func NewApplication(lc fx.Lifecycle, conf *Config) *echo.Echo {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(m.Zerolog)

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Start(conf.Settings.Server.Port); err != nil && err != http.ErrServerClosed {
					log.Fatal().Err(err).Msg("failed to start server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	return e
}
