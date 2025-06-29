package home

import (
	"github.com/labstack/echo/v4"

	"tidebot/pkg/common"
	"tidebot/pkg/middleware"
	"tidebot/pkg/ui/layout"
)

func RegisterHomeRoutes(e *echo.Echo) {
	e.GET("/", middleware.HTMLRoute(
		func(c echo.Context) error {
			req := c.Request()
			userAgent := req.Header.Get("User-Agent")
			isMobile := common.IsMobileUserAgent(userAgent)

			return layout.RenderPage(
				c,
				200,
				HomePage(isMobile),
			)
		},
	))
}
