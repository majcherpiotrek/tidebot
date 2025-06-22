package whatsapp

import "github.com/labstack/echo/v4"

func RegisterWhatsappWebhook(e *echo.Echo) {

	e.POST("/message", func(c echo.Context) error {
		return nil
	})
}
