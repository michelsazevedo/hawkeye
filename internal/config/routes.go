package config

import (
	"github.com/labstack/echo/v4"
	"github.com/michelsazevedo/hawkeye/internal/api"
	"github.com/michelsazevedo/hawkeye/internal/domain"
)

func Routes(e *echo.Echo, courseHandler api.SearchHandler[domain.Course]) {
	courses := e.Group("/", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Get("subdomain") == "courses" {
				return next(c)
			}
			return echo.ErrNotFound
		}
	})

	courses.GET("search", courseHandler.Search)
}
