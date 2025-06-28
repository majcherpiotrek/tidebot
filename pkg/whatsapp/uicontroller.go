package whatsapp

import (
	"github.com/labstack/echo/v4"

	"tidebot/pkg/whatsapp/qrcode"
	"tidebot/pkg/whatsapp/whatsappsignup"
)

func RegisterWhatsappComponents(e *echo.Echo) {
	qrcode.RegisterComponent(e)
	whatsappsignup.RegisterComponent(e)
}
