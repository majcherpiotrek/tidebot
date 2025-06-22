package whatsapp

import (
	"fmt"
	"strings"
	"tidebot/pkg/users/services"
	"tidebot/pkg/worldtides"

	"github.com/labstack/echo/v4"
)

type WhatsAppService interface {
	ProcessMessage(body string, from string, profileName *string) error
	SendTideExtremesMessage(phoneNumber string, extremes []worldtides.Extreme, date string) error
}

type whatsappServiceImpl struct {
	userService    services.UserService
	whatsappClient WhatsappClient
	log            echo.Logger
}

func NewWhatsAppService(userService services.UserService, whatsappClient WhatsappClient, log echo.Logger) WhatsAppService {
	return &whatsappServiceImpl{
		userService:    userService,
		whatsappClient: whatsappClient,
		log:            log,
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

func (s *whatsappServiceImpl) SendTideExtremesMessage(phoneNumber string, extremes []worldtides.Extreme, date string) error {
	s.log.Debugf("Sending tide extremes message to %s for date %s", phoneNumber, date)

	message := s.formatTideExtremesMessage(extremes, date)

	err := s.whatsappClient.SendMessage(message, phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to send tide extremes message to %s: %w", phoneNumber, err)
	}

	s.log.Infof("Successfully sent tide extremes message to %s", phoneNumber)
	return nil
}

func (s *whatsappServiceImpl) formatTideExtremesMessage(extremes []worldtides.Extreme, date string) string {
	if len(extremes) == 0 {
		return fmt.Sprintf("🌊 *Tide Report for %s*\n\nNo tide data available for today.", date)
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("🌊 *Tide Report for %s*\n\n", date))

	for _, extreme := range extremes {
		tideTime := extreme.Time().Format("15:04")
		var emoji string
		var extraNewLine string
		if extreme.IsHighTide() {
			emoji = "⬆️"
			extraNewLine = "\n"
		} else {
			emoji = "⬇️"
			extraNewLine = ""
		}

		message.WriteString(fmt.Sprintf("%s *%s Tide*: %s (%.2fm)%s\n",
			emoji, extreme.Type, tideTime, extreme.Height, extraNewLine))
	}

	message.WriteString("\n📍 Fuerteventura, Risco del Paso, Canary Islands")

	return message.String()
}
