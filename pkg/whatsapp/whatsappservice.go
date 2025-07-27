package whatsapp

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"tidebot/pkg/common"
	"tidebot/pkg/environment"
	"tidebot/pkg/notifications/repositories"
	"tidebot/pkg/users/services"
	"tidebot/pkg/worldtides"
	"time"

	"github.com/labstack/echo/v4"
)

type WhatsAppService interface {
	ProcessMessage(body string, from string, profileName *string) error
	SendTideExtremesMessage(phoneNumber string, extremes []worldtides.Extreme, date time.Time) error
	SendDailyTideNotification(phoneNumber string, userName string, extremes []worldtides.Extreme) error
}

type whatsappServiceImpl struct {
	userService                        services.UserService
	notificationSubscriptionRepository repositories.NotificationSubscriptionRepository
	worldTidesClient                   worldtides.WorldTidesClient
	whatsappClient                     WhatsappClient
	log                                echo.Logger
}

func NewWhatsAppService(userService services.UserService, notificationSubscriptionRepository repositories.NotificationSubscriptionRepository, worldTidesClient worldtides.WorldTidesClient, whatsappClient WhatsappClient, log echo.Logger) WhatsAppService {
	return &whatsappServiceImpl{
		userService:                        userService,
		notificationSubscriptionRepository: notificationSubscriptionRepository,
		worldTidesClient:                   worldTidesClient,
		whatsappClient:                     whatsappClient,
		log:                                log,
	}
}

func (s *whatsappServiceImpl) ProcessMessage(body string, from string, profileName *string) error {
	s.log.Debugf("Processing WhatsApp message - body: %s, from: %s, profileName: %v", body, from, profileName)

	cleanPhoneNumber := strings.TrimPrefix(from, "whatsapp:")

	trimmedBody := strings.ToLower(strings.TrimSpace(body))

	whitespaceRegexp := regexp.MustCompile(`\s+`)
	commandWithArguments := whitespaceRegexp.Split(trimmedBody, -1)

	if len(commandWithArguments) < 1 {
		return s.defaultMessageHandler(cleanPhoneNumber, profileName)
	}

	command := commandWithArguments[0]
	arguments := commandWithArguments[1:]

	switch command {
	case "tides":
		return s.handleTidesCommand(cleanPhoneNumber, arguments)
	case "start":
		return s.handleStartCommand(cleanPhoneNumber, profileName)
	case "stop":
		return s.handleStopCommand(cleanPhoneNumber)
	default:
		return s.defaultMessageHandler(cleanPhoneNumber, profileName)

	}
}

func (s *whatsappServiceImpl) SendTideExtremesMessage(phoneNumber string, extremes []worldtides.Extreme, date time.Time) error {
	s.log.Debugf("Sending tide extremes message to %s for date %s", phoneNumber, date)

	message := s.formatTideExtremesMessage(extremes, date)

	err := s.whatsappClient.SendMessage(message, phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to send tide extremes message to %s: %w", phoneNumber, err)
	}

	s.log.Infof("Successfully sent tide extremes message to %s", phoneNumber)
	return nil
}

func (s *whatsappServiceImpl) formatTideExtremesMessage(extremes []worldtides.Extreme, date time.Time) string {
	dateFormatted := date.Format("Monday, 2006-01-02")
	if len(extremes) == 0 {
		return fmt.Sprintf("🌊 *Tides for %s*\n\nNo tide data available for today.", dateFormatted)
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintf("🌊 *Tides for %s*\n\n", dateFormatted))

	// Load Atlantic/Canary timezone
	canaryTZ, err := time.LoadLocation("Atlantic/Canary")
	if err != nil {
		s.log.Errorf("Failed to load Atlantic/Canary timezone: %v", err)
		canaryTZ = time.UTC // Fallback to UTC
	}

	for _, extreme := range extremes {
		// Convert time to Canary timezone
		tideTimeInCanary := extreme.Time().In(canaryTZ)
		tideTime := tideTimeInCanary.Format("15:04")

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

func (s *whatsappServiceImpl) defaultMessageHandler(phoneNumber string, profileName *string) error {
	s.log.Info("Received message, saving user")

	_, err := s.userService.GetUserByPhoneNumber(phoneNumber)
	isNewUser := false

	if err != nil {
		_, err = s.userService.SaveUser(phoneNumber, profileName)

		if err != nil {
			return fmt.Errorf("Failed to save user with phone number %s: %w", phoneNumber, err)
		}

		isNewUser = true
	}

	err = s.sendWelcomeMessage(phoneNumber, profileName, isNewUser)
	if err != nil {
		s.log.Errorf("Failed to send welcome message to %s: %v", phoneNumber, err)
	}

	s.log.Info("Successfully processed message")
	return nil
}

type TidesResponseForDay struct {
	Day           time.Time
	TidesResponse *worldtides.WorldTidesResponse
	Err           error
}

func (s *whatsappServiceImpl) handleTidesCommand(phoneNumber string, arguments []string) error {
	s.log.Infof("Handling tides command for %s. Arguments: %v", phoneNumber, arguments)

	var dates []time.Time

	if len(arguments) > 0 {
		dates = s.parseTidesCommandArguments(arguments)
	} else {
		dates = append(dates, common.Today())
	}

	ch := make(chan TidesResponseForDay, len(dates))

	for _, day := range dates {
		go s.getTidesWorker(day, ch)
	}

	var responses []TidesResponseForDay
	for range dates {
		res := <-ch
		responses = append(responses, res)
	}

	sort.Slice(responses, func(i, j int) bool {
		return responses[i].Day.Before(responses[j].Day)
	})

	for _, response := range responses {
		if response.Err != nil {
			s.whatsappClient.SendMessage(fmt.Sprintf("❌ Sorry, I couldn't fetch tide data for %s. Please try again later.", response.Day.Format("2006-01-02")), phoneNumber)
		} else {
			s.SendTideExtremesMessage(phoneNumber, response.TidesResponse.Extremes, response.Day)
		}
	}

	return nil
}

func (s *whatsappServiceImpl) getTidesWorker(day time.Time, results chan<- TidesResponseForDay) {
	tidesResponse, err := s.worldTidesClient.GetTides(day)
	dayFormatted := day.Format("2006-01-02")

	if err != nil {
		wrappedError := fmt.Errorf("Failed to fetch tide extremes for day %s: %v", dayFormatted, err)
		s.log.Errorf("%v", wrappedError)
		results <- TidesResponseForDay{
			Day:           day,
			TidesResponse: nil,
			Err:           wrappedError,
		}
	} else {
		results <- TidesResponseForDay{
			Day:           day,
			TidesResponse: tidesResponse,
			Err:           nil,
		}
	}
}

func (s *whatsappServiceImpl) parseTidesCommandArguments(args []string) []time.Time {
	argsClean := slices.Clone(args)

	for i := range argsClean {
		argsClean[i] = strings.ToLower(argsClean[i])
	}

	containsWeekArg := slices.Contains(argsClean, "week")

	if containsWeekArg {
		week := make([]time.Time, 7)

		for i := range week {
			week[i] = common.Today().Add(time.Duration(i*24) * time.Hour)
		}

		return week
	}

	var dates []time.Time

	for i := range argsClean {
		date, err := common.ParseDate(argsClean[i])

		if err != nil {
			continue
		}

		dates = append(dates, date)
	}

	return dates
}

func (s *whatsappServiceImpl) handleStartCommand(phoneNumber string, profileName *string) error {
	s.log.Infof("Handling start command for %s", phoneNumber)

	user, err := s.userService.SaveUser(phoneNumber, profileName)
	if err != nil {
		s.log.Errorf("Failed to save user for phone %s: %v", phoneNumber, err)
		return s.whatsappClient.SendMessage("❌ Sorry, there was an error. Please try again later.", phoneNumber)
	}

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

	user, err := s.userService.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		s.log.Warnf("User not found for phone %s, cannot stop notifications", phoneNumber)
		return s.whatsappClient.SendMessage("🤷‍♂️ You don't have any active notifications to stop.\n\nSend *start* to enable tide notifications!", phoneNumber)
	}

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

func (s *whatsappServiceImpl) sendWelcomeMessage(phoneNumber string, profileName *string, isNewUser bool) error {
	personalizedWelcome := "Hi!"

	if profileName != nil {
		personalizedWelcome = fmt.Sprintf(`Hi %s!`, *profileName)
	}

	newUserMessage := " Welcome to TideBot!"
	if !isNewUser {
		newUserMessage = ""
	}

	welcomeMessage := fmt.Sprintf(`🌊 *%s%s*

Tide reports for *Risco del Paso, Fuerteventura*.

Your tide reports include high and low tide times with precise heights 🏄‍♂️

%s
`, personalizedWelcome, newUserMessage, AVAILABLE_COMMANDS)

	err := s.whatsappClient.SendMessage(welcomeMessage, phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to send welcome message: %w", err)
	}

	s.sendQuickReplyMessage(phoneNumber)

	s.log.Infof("Sent welcome message to %s", phoneNumber)
	return nil
}

func (s *whatsappServiceImpl) sendQuickReplyMessage(phoneNumber string) {
	s.log.Infof("Attempting to send interactive template %s to %s", QUICK_REPLY_MESSAGE_TEMPLATE_SID, phoneNumber)
	err := s.whatsappClient.SendInteractiveTemplate(QUICK_REPLY_MESSAGE_TEMPLATE_SID, phoneNumber)
	if err != nil {
		s.log.Errorf("Failed to send interactive template to %s: %v", phoneNumber, err)
	} else {
		s.log.Infof("Successfully sent interactive template to %s", phoneNumber)
	}
}

func (s *whatsappServiceImpl) SendDailyTideNotification(phoneNumber string, userName string, extremes []worldtides.Extreme) error {
	// Check environment - use text message in development, template in production
	env := os.Getenv("GO_ENV")

	variables := s.buildDailyTidesNotificationVariables(userName, extremes)

	if env == string(environment.EnvDevelopment) {
		s.log.Infof("Using text message for daily notification in development environment")
		return s.sendDailyTideNotificationAsText(phoneNumber, variables)
	}

	s.log.Infof("Sending daily tide notification template to %s", phoneNumber)

	if len(extremes) < 4 {
		return fmt.Errorf("insufficient tide extremes: need 4, got %d", len(extremes))
	}

	err := s.whatsappClient.SendTemplateWithVariables(DAILY_TIDE_NOTIFICATION_TEMPLATE_SID, variables, phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to send daily tide notification: %w", err)
	}

	s.log.Infof("Successfully sent daily tide notification to %s", phoneNumber)
	return nil
}

func (s *whatsappServiceImpl) sendDailyTideNotificationAsText(phoneNumber string, variables []string) error {
	s.log.Infof("Sending daily tide notification as text to %s", phoneNumber)

	if len(variables) != 9 {
		return fmt.Errorf("9 variables required, got %d", len(variables))
	}

	// Build the message text
	var message strings.Builder
	message.WriteString(fmt.Sprintf("Hi %s!\n\n", variables[8]))
	message.WriteString("Here is your daily tide report:\n\n")

	// Process each extreme
	for i, j := 0, 1; i < 7; i, j = i+2, j+1 {
		tideType := variables[i]
		tideInfo := variables[i+1]

		message.WriteString(fmt.Sprintf("  %d. %s tide: %s\n", j, tideType, tideInfo))
	}

	message.WriteString("\nLocation: Risco del Paso, Fuerteventura\n\n")
	message.WriteString("If you don't want to receive those notifications anymore, reply 'stop' to this message. Have a great day on the water!")

	err := s.whatsappClient.SendMessage(message.String(), phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to send daily tide notification as text: %w", err)
	}

	s.log.Infof("Successfully sent daily tide notification as text to %s", phoneNumber)
	return nil
}

func (s *whatsappServiceImpl) buildDailyTidesNotificationVariables(userName string, extremes []worldtides.Extreme) []string {
	if len(extremes) < 3 {
		return []string{}
	}

	// TODO: correlate timezone with place when adding support for more location
	tz := s.getTimezone()

	nowInCanary := time.Now().In(tz)
	todayInCanary := nowInCanary.Format("2006-01-02")

	// Build variables for the template
	variables := make([]string, 9) // 9 variables total

	for i := range variables {
		variables[i] = ""
	}

	// {{9}} - User name
	variables[8] = userName

	// Process each extreme ({{1}} through {{8}})
	for i := range extremes {
		extreme := extremes[i]

		// Convert time to Canary timezone
		tideTimeInCanary := extreme.Time().In(tz)

		// Determine if this tide is on the next day (compared to Canary time)
		daySuffix := ""
		if tideTimeInCanary.Format("2006-01-02") != todayInCanary {
			daySuffix = " (+1 day)"
		}

		// {{1}} - index 0, {{3}} - index 2, {{5}} - index 4, {{7}} - index 6 -- Tide type (High/Low)
		if extreme.IsHighTide() {
			variables[i*2] = "High"
		} else {
			variables[i*2] = "Low"
		}

		// {{2}} - index 1, {{4}} - index 3, {{6}} - index 5, {{8}} - index 7 -- Time and height in Canary timezone
		tideTime := tideTimeInCanary.Format("15:04")
		variables[i*2+1] = fmt.Sprintf("%s%s (%.2fm)", tideTime, daySuffix, extreme.Height)
	}

	return variables
}

func (s *whatsappServiceImpl) getTimezone() *time.Location {
	canaryTZ, err := time.LoadLocation("Atlantic/Canary")
	if err != nil {
		s.log.Errorf("Failed to load Atlantic/Canary timezone: %v", err)
		canaryTZ = time.UTC // Fallback to UTC
	}

	return canaryTZ
}

const AVAILABLE_COMMANDS = `
*Available commands:*
📱 Send *tides* - Get today's tide info
   Examples: _tides tomorrow_, _tides week_, _tides today tomorrow_, _tides today 24/12/2025_
🔔 Send *start* - Enable daily notifications  
🔕 Send *stop* - Disable notifications
`

const QUICK_REPLY_MESSAGE_TEMPLATE_SID = "HX6f156e3466407a835bef6505f85cf9b1"
const DAILY_TIDE_NOTIFICATION_TEMPLATE_SID = "HX7161523078d66056973776cbf70f583a"
