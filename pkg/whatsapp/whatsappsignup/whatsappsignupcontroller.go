package whatsappsignup

import (
	"strings"

	"github.com/labstack/echo/v4"

	"tidebot/pkg/ui/layout"
)

func RegisterComponent(e *echo.Echo) {
	e.GET("/whatsapp/components/whatsappsignup",
		func(c echo.Context) error {
			req := c.Request()

			userAgent := req.Header.Get("User-Agent")
			isMobile := isMobileUserAgent(userAgent)

			return layout.RenderComponent(
				c,
				200,
				WhatsAppSignUp(isMobile, "+14155238886", "overpowered"),
			)
		},
	)
}

func isMobileUserAgent(userAgent string) bool {
	mobileAgents := []string{
		"Android", "iPhone", "iPad", "iPod", "Mobile",
		"BlackBerry", "Windows Phone", "Opera Mini",
	}

	userAgent = strings.ToLower(userAgent)
	for _, agent := range mobileAgents {
		if strings.Contains(userAgent, strings.ToLower(agent)) {
			return true
		}
	}
	return false
}
