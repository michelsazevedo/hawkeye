package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func Subdomain(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		host := c.Request().Host
		schema := strings.Split(host, ".")

		if len(schema) > 2 {
			c.Set("subdomain", schema[0])
			return next(c)
		}
		return echo.ErrNotFound
	}
}
