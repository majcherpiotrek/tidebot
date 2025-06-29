package whatsapp

import (
	"fmt"
	"strings"
	"time"
	"tidebot/pkg/notifications/repositories"
	"tidebot/pkg/users/services"
	"tidebot/pkg/worldtides"

	"github.com/labstack/echo/v4"
)

type WhatsAppService interface {
	ProcessMessage(body string, from string, profileName *string) error
	SendTideExtremesMessage(phoneNumber string, extremes []worldtides.Extreme, date string) error
}

type whatsappServiceImpl struct {
	userService                      services.UserService
	notificationSubscriptionRepository repositories.NotificationSubscriptionRepository
	worldTidesClient                 worldtides.WorldTidesClient
	whatsappClient                   WhatsappClient
	log                              echo.Logger
}

func NewWhatsAppService(userService services.UserService, notificationSubscriptionRepository repositories.NotificationSubscriptionRepository, worldTidesClient worldtides.WorldTidesClient, whatsappClient WhatsappClient, log echo.Logger) WhatsAppService {
	return &whatsappServiceImpl{
		userService:                      userService,
		notificationSubscriptionRepository: notificationSubscriptionRepository,
		worldTidesClient:                 worldTidesClient,
		whatsappClient:                   whatsappClient,
		log:                              log,
	}
}

func (s *whatsappServiceImpl) ProcessMessage(body string, from string, profileName *string) error {
	s.log.Debugf("Processing WhatsApp message - body: %s, from: %s, profileName: %v", body, from, profileName)

	// Handle commands
	command := strings.ToLower(strings.TrimSpace(body))
	cleanPhoneNumber := strings.TrimPrefix(from, "whatsapp:")
	
	switch command {
	case "overpowered":
		return s.handleSignupMessage(cleanPhoneNumber, profileName)
	case "tides":
		return s.handleTidesCommand(cleanPhoneNumber)
	case "start":
		return s.handleStartCommand(cleanPhoneNumber)
	case "stop":
		return s.handleStopCommand(cleanPhoneNumber)
	default:
		s.log.Debugf("Message body '%s' does not match any known commands, ignoring", body)
		return nil
	}
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

func (s *whatsappServiceImpl) handleSignupMessage(phoneNumber string, profileName *string) error {
	s.log.Info("Received 'overpowered' message, saving user")

	err := s.userService.SaveUser(phoneNumber, profileName)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	// Send welcome message
	err = s.sendWelcomeMessage(phoneNumber)
	if err != nil {
		s.log.Errorf("Failed to send welcome message to %s: %v", phoneNumber, err)
		// Don't return error - user is saved, just welcome message failed
	}

	s.log.Info("Successfully processed 'overpowered' message")
	return nil
}

func (s *whatsappServiceImpl) handleTidesCommand(phoneNumber string) error {
	s.log.Infof("Handling tides command for %s", phoneNumber)

	// Get today's date
	today := time.Now().Format("2006-01-02")
	
	// Fetch tide extremes from WorldTides API
	tidesResponse, err := s.worldTidesClient.GetTidalExtremesForDay(today)
	if err != nil {
		s.log.Errorf("Failed to fetch tide extremes for %s: %v", phoneNumber, err)
		return s.whatsappClient.SendMessage("❌ Sorry, I couldn't fetch tide data right now. Please try again later.", phoneNumber)
	}

	// Send tide extremes using existing method
	return s.SendTideExtremesMessage(phoneNumber, tidesResponse.Extremes, today)
}

func (s *whatsappServiceImpl) handleStartCommand(phoneNumber string) error {
	s.log.Infof("Handling start command for %s", phoneNumber)

	// Auto-register user if not exists (idempotent)
	err := s.userService.SaveUser(phoneNumber, nil)
	if err != nil {
		s.log.Errorf("Failed to save user for phone %s: %v", phoneNumber, err)
		return s.whatsappClient.SendMessage("❌ Sorry, there was an error. Please try again later.", phoneNumber)
	}

	// Get user to get ID
	user, err := s.userService.GetUserByPhoneNumber(phoneNumber)
	if err != nil || user == nil {
		s.log.Errorf("Failed to get user after save for phone %s: %v", phoneNumber, err)
		return s.whatsappClient.SendMessage("❌ Sorry, there was an error. Please try again later.", phoneNumber)
	}

	// Create/enable subscription
	err = s.notificationSubscriptionRepository.CreateSubscription(user.ID)
	if err != nil {
		s.log.Errorf("Failed to create subscription for user %d: %v", user.ID, err)
		return s.whatsappClient.SendMessage("❌ Sorry, there was an error enabling notifications. Please try again later.", phoneNumber)
	}

	confirmationMessage := `🔔 *Notifications Enabled!*

You'll now receive daily tide reports for *Risco del Paso, Fuerteventura* every morning.

📱 Send *tides* anytime for current tide info
🔕 Send *stop* to disable notifications

Welcome aboard! 🌊`

	return s.whatsappClient.SendMessage(confirmationMessage, phoneNumber)
}

func (s *whatsappServiceImpl) handleStopCommand(phoneNumber string) error {
	s.log.Infof("Handling stop command for %s", phoneNumber)

	// Get user
	user, err := s.userService.GetUserByPhoneNumber(phoneNumber)
	if err != nil || user == nil {
		s.log.Warnf("User not found for phone %s, cannot stop notifications", phoneNumber)
		return s.whatsappClient.SendMessage("🤷‍♂️ You don't have any active notifications to stop.\n\nSend *start* to enable tide notifications!", phoneNumber)
	}

	// Disable subscription
	err = s.notificationSubscriptionRepository.DisableSubscription(user.ID)
	if err != nil {
		if strings.Contains(err.Error(), "no subscription found") {
			return s.whatsappClient.SendMessage("🤷‍♂️ You don't have any active notifications to stop.\n\nSend *start* to enable tide notifications!", phoneNumber)
		}
		s.log.Errorf("Failed to disable subscription for user %d: %v", user.ID, err)
		return s.whatsappClient.SendMessage("❌ Sorry, there was an error. Please try again later.", phoneNumber)
	}

	confirmationMessage := `🔕 *Notifications Disabled*

You'll no longer receive daily tide reports.

📱 Send *tides* anytime for current tide info
🔔 Send *start* to re-enable notifications

Thanks for using TideBot! 🌊`

	return s.whatsappClient.SendMessage(confirmationMessage, phoneNumber)
}

func (s *whatsappServiceImpl) sendWelcomeMessage(phoneNumber string) error {
	welcomeMessage := `🌊 *Welcome to TideBot!*

Great! You're now registered to receive tide reports for *Risco del Paso, Fuerteventura*.

Your tide reports include high and low tide times with precise heights. Perfect for planning your beach day, surfing, or fishing! 🏄‍♂️🎣

*Available commands:*
📱 Send *tides* - Get current tide info
🔔 Send *start* - Enable daily notifications  
🔕 Send *stop* - Disable notifications`

	// Send welcome message first
	err := s.whatsappClient.SendMessage(welcomeMessage, phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to send welcome message: %w", err)
	}

	// Send interactive template as follow-up
	templateSID := "HX6f156e3466407a835bef6505f85cf9b1"
	err = s.whatsappClient.SendInteractiveTemplate(templateSID, phoneNumber)
	if err != nil {
		s.log.Warnf("Failed to send interactive template to %s: %v", phoneNumber, err)
		// Don't return error - welcome message was sent successfully
	}

	s.log.Infof("Sent welcome message to %s", phoneNumber)
	return nil
}
