package whatsapp

import (
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func RegisterWhatsappWebhook(e *echo.Echo, whatsappService WhatsAppService) {

	e.POST("/message", func(c echo.Context) error {
		logger := c.Echo().Logger

		logger.Infof("📱 WhatsApp webhook received from IP: %s", c.RealIP())
		logger.Infof("📱 Request: %s %s", c.Request().Method, c.Request().URL.Path)

		if len(c.QueryParams()) > 0 {
			logger.Info("📱 Query parameters:")
			for key, values := range c.QueryParams() {
				for _, value := range values {
					logger.Infof("  %s: %s", key, value)
				}
			}
		}

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			logger.Errorf("📱 Failed to read request body: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Failed to read request body",
			})
		}

		logger.Infof("📱 Raw body (%d bytes): %s", len(body), string(body))

		// Parse form data (Twilio sends form-encoded data)
		formData, err := url.ParseQuery(string(body))
		if err != nil {
			logger.Errorf("📱 Failed to parse form data: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Failed to parse form data",
			})
		} else if len(formData) > 0 {
			logger.Info("📱 Form data:")
			for key, values := range formData {
				for _, value := range values {
					logger.Infof("  %s: %s", key, value)
				}
			}

			// Extract message details
			messageBody := formData.Get("Body")
			from := formData.Get("From")
			profileName := formData.Get("ProfileName")

			// Process the message if we have the required fields
			if messageBody != "" && from != "" {
				var profileNamePtr *string
				if profileName != "" {
					profileNamePtr = &profileName
				}

				err := whatsappService.ProcessMessage(messageBody, from, profileNamePtr)
				if err != nil {
					logger.Errorf("📱 Failed to process message: %v", err)
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "Failed to process message",
					})
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
