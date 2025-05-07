package config

import (
	"github.com/labstack/echo/v4"
	"github.com/michelsazevedo/hawkeye/api"
	"github.com/michelsazevedo/hawkeye/domain"
)

func Routes(e *echo.Echo, courseHandler api.SearchHandler[domain.Course]) {
	courses := e.Group("courses")
	courses.GET("/search", courseHandler.Search)
}
