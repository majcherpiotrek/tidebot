package whatsappsignup

import (
	"github.com/labstack/echo/v4"

	"tidebot/pkg/common"
	"tidebot/pkg/ui/layout"
)

func RegisterComponent(e *echo.Echo) {
	e.GET(whatsAppSignupComponentURI,
		func(c echo.Context) error {
			req := c.Request()

			userAgent := req.Header.Get("User-Agent")
			isMobile := common.IsMobileUserAgent(userAgent)

			phoneNumber := "+14155238886"

			return layout.RenderComponent(
				c,
				200,
				whatsAppSignUpContent(WhatsAppSignUpProps{
					IsMobile:    isMobile,
					PhoneNumber: phoneNumber,
					Message:     "overpowered",
				}),
			)
		},
	)
}
