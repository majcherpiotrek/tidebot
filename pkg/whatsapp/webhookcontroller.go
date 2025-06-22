package whatsapp

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterWhatsappWebhook(e *echo.Echo) {

	e.POST("/message", func(c echo.Context) error {
		logger := c.Echo().Logger

		// Log basic request info
		logger.Infof("📱 WhatsApp webhook received from IP: %s", c.RealIP())
		logger.Infof("📱 Request method: %s %s", c.Request().Method, c.Request().URL.Path)

		// Log all headers
		logger.Info("📱 Headers:")
		for name, values := range c.Request().Header {
			for _, value := range values {
				logger.Infof("  %s: %s", name, value)
			}
		}

		// Log query parameters
		if len(c.QueryParams()) > 0 {
			logger.Info("📱 Query parameters:")
			for key, values := range c.QueryParams() {
				for _, value := range values {
					logger.Infof("  %s: %s", key, value)
				}
			}
		}

		// Read and log the raw body
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			logger.Errorf("📱 Failed to read request body: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Failed to read request body",
			})
		}

		logger.Infof("📱 Raw body (%d bytes): %s", len(body), string(body))

		// Log form data (Twilio sends form-encoded data)
		if err := c.Request().ParseForm(); err == nil && len(c.Request().Form) > 0 {
			logger.Info("📱 Form data:")
			for key, values := range c.Request().Form {
				for _, value := range values {
					logger.Infof("  %s: %s", key, value)
				}
			}
		}

		// Respond with success (Twilio expects 200 OK)
		logger.Info("📱 Responding with 200 OK")
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "received",
			"message": "Webhook processed successfully",
		})
	})

	// Optional: Add a simple GET endpoint for testing
	e.GET("/message", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "webhook_ready",
			"message": "WhatsApp webhook endpoint is ready",
			"method":  "POST /message for webhooks",
		})
	})
}
