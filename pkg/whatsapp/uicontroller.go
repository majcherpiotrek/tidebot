package whatsapp

import (
	"github.com/labstack/echo/v4"

	"tidebot/pkg/whatsapp/whatsappsignup"
)

func RegisterComponents(e *echo.Echo) {
	whatsappsignup.RegisterComponent(e)
}
