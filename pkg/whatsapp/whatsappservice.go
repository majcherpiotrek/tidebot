package whatsapp

import (
	"fmt"
	"strings"
	"tidebot/pkg/users/services"

	"github.com/labstack/echo/v4"
)

type WhatsAppService interface {
	ProcessMessage(body string, from string, profileName *string) error
}

type whatsappServiceImpl struct {
	userService services.UserService
	log         echo.Logger
}

func NewWhatsAppService(userService services.UserService, log echo.Logger) WhatsAppService {
	return &whatsappServiceImpl{
		userService: userService,
		log:         log,
	}
}

func (s *whatsappServiceImpl) ProcessMessage(body string, from string, profileName *string) error {
	s.log.Debugf("Processing WhatsApp message - body: %s, from: %s, profileName: %v", body, from, profileName)

	// Check if message body is "overpowered"
	if strings.ToLower(strings.TrimSpace(body)) == "overpowered" {
		s.log.Info("Received 'overpowered' message, saving user")

		cleanPhoneNumber := strings.TrimPrefix(from, "whatsapp:")

		err := s.userService.SaveUser(cleanPhoneNumber, profileName)
		if err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}

		s.log.Info("Successfully processed 'overpowered' message")
		return nil
	}

	s.log.Debugf("Message body '%s' does not match 'overpowered', ignoring", body)
	return nil
}

