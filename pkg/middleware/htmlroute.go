package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func HTMLRoute(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !IsHTMLRequest(c) {
			return c.JSON(http.StatusNotAcceptable, map[string]string{
				"error": "This endpoint only serves HTML",
			})
		}
		return handler(c)
	}
}

func IsHTMLRequest(c echo.Context) bool {
	accept := c.Request().Header.Get("Accept")
	return accept == "" ||
		strings.Contains(accept, "text/html") ||
		strings.Contains(accept, "application/xhtml+xml")
}
