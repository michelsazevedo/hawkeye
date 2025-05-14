package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func Zerolog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)
		stop := time.Now()

		log.Info().
			Str("Method=", c.Request().Method).
			Str("URI", c.Request().RequestURI).
			Str("RemoteIp", c.RealIP()).
			Str("UserAgent", c.Request().UserAgent()).
			Int("Sstatus", c.Response().Status).
			Dur("latency", stop.Sub(start)).
			Msg("Host=" + c.Request().Host)

		return err
	}
}
