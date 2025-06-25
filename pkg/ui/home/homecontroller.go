package home

import (
	"github.com/labstack/echo/v4"

	"tidebot/pkg/middleware"
	"tidebot/pkg/ui/layout"
)

func RegisterHomeRoutes(e *echo.Echo) {
	e.GET("/", middleware.HTMLRoute(
		func(c echo.Context) error {
			return layout.RenderPage(
				c,
				200,
				HomePage(),
			)
		},
	))
}
