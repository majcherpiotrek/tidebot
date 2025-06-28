package qrcode

import (
	"github.com/labstack/echo/v4"

	"tidebot/pkg/ui/layout"
)

func RegisterComponent(e *echo.Echo) {
	e.GET("/whatsapp/components/qrcode",
		func(c echo.Context) error {
			return layout.RenderComponent(
				c,
				200,
				QrCode("+14155238886", "overpowered"),
			)
		},
	)
}
