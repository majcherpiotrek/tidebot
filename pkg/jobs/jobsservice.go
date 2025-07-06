package jobs

import (
	"fmt"
	"tidebot/pkg/notifications/repositories"
	"tidebot/pkg/users/services"
	"tidebot/pkg/whatsapp"
	"tidebot/pkg/worldtides"
	"time"

	"github.com/labstack/echo/v4"
)

type JobsService interface {
	SendTideExtremesToAllUsers() error
	SendDailyNotificationsV2() (int, error)
}

type jobsServiceImpl struct {
	userService                        services.UserService
	notificationSubscriptionRepository repositories.NotificationSubscriptionRepository
	whatsappService                    whatsapp.WhatsAppService
	worldTidesClient                   worldtides.WorldTidesClient
	log                                echo.Logger
}

func NewJobsService(
	userService services.UserService,
	notificationSubscriptionRepository repositories.NotificationSubscriptionRepository,
	whatsappService whatsapp.WhatsAppService,
	worldTidesClient worldtides.WorldTidesClient,
	log echo.Logger,
) JobsService {
	return &jobsServiceImpl{
		userService:                        userService,
		notificationSubscriptionRepository: notificationSubscriptionRepository,
		whatsappService:                    whatsappService,
		worldTidesClient:                   worldTidesClient,
		log:                                log,
	}
}

func (j *jobsServiceImpl) SendTideExtremesToAllUsers() error {
	j.log.Info("Starting job: Send tide extremes to all users")

	// Get today's date
	today := time.Now().Format("2006-01-02")
	j.log.Debugf("Fetching tide extremes for date: %s", today)

	// Fetch tide extremes from WorldTides API
	tidesResponse, err := j.worldTidesClient.GetTides(today)
	if err != nil {
		return fmt.Errorf("failed to fetch tide extremes: %w", err)
	}

	j.log.Debugf("Received %d tide extremes for %s", len(tidesResponse.Extremes), today)

	// Debug log all extremes received
	for i, extreme := range tidesResponse.Extremes {
		localTime := extreme.Time().Format("2006-01-02 15:04:05 MST")
		j.log.Debugf("Extreme %d: %s tide at %s (%.4fm) - Unix: %d, Date field: %s",
			i+1, extreme.Type, localTime, extreme.Height, extreme.Dt, extreme.Date)
	}

	// Get all users with enabled subscriptions
	subscriptions, err := j.notificationSubscriptionRepository.GetEnabledSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to get enabled subscriptions: %w", err)
	}

	j.log.Infof("Found %d users with enabled subscriptions to send tide extremes", len(subscriptions))

	// Send tide extremes to each subscribed user
	successCount := 0
	errorCount := 0

	for _, subscription := range subscriptions {
		// Get user details to get phone number
		user, err := j.userService.GetUserByID(subscription.UserID)
		if err != nil {
			j.log.Errorf("Failed to get user details for subscription ID=%d, UserID=%d: %v", subscription.ID, subscription.UserID, err)
			errorCount++
			continue
		}

		j.log.Debugf("Sending tide extremes to subscribed user ID=%d, phone=%s", subscription.UserID, user.PhoneNumber)

		err = j.whatsappService.SendTideExtremesMessage(user.PhoneNumber, tidesResponse.Extremes, today)
		if err != nil {
			j.log.Errorf("Failed to send tide extremes to user ID=%d: %v", subscription.UserID, err)
			errorCount++
		} else {
			j.log.Debugf("Successfully sent tide extremes to user ID=%d", subscription.UserID)
			successCount++
		}
	}

	j.log.Infof("Job completed: %d successful, %d errors out of %d subscribed users", successCount, errorCount, len(subscriptions))

	if errorCount > 0 {
		return fmt.Errorf("job completed with %d errors out of %d subscribed users", errorCount, len(subscriptions))
	}

	return nil
}

func (j *jobsServiceImpl) SendDailyNotificationsV2() (int, error) {
	j.log.Info("Starting job: Send daily tide notifications (v2)")

	// Get today's date
	today := time.Now().Format("2006-01-02")
	j.log.Debugf("Fetching tide extremes for date: %s", today)

	// Fetch tide extremes from WorldTides API
	tidesResponse, err := j.worldTidesClient.GetTides(today)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch tide extremes: %w", err)
	}

	j.log.Debugf("Received %d tide extremes for %s", len(tidesResponse.Extremes), today)

	// Ensure we have enough extremes for the notification
	if len(tidesResponse.Extremes) < 4 {
		return 0, fmt.Errorf("insufficient tide extremes for daily notification: need 4, got %d", len(tidesResponse.Extremes))
	}

	// Get all users with enabled subscriptions
	subscriptions, err := j.notificationSubscriptionRepository.GetEnabledSubscriptions()
	if err != nil {
		return 0, fmt.Errorf("failed to get enabled subscriptions: %w", err)
	}

	j.log.Infof("Found %d users with enabled subscriptions to send daily notifications", len(subscriptions))

	// Send daily notifications to each subscribed user
	successCount := 0
	errorCount := 0

	for _, subscription := range subscriptions {
		// Get user details to get phone number and name
		user, err := j.userService.GetUserByID(subscription.UserID)
		if err != nil {
			j.log.Errorf("Failed to get user details for subscription ID=%d, UserID=%d: %v", subscription.ID, subscription.UserID, err)
			errorCount++
			continue
		}

		// Use name if available, otherwise use phone number
		userName := user.PhoneNumber
		if user.Name != nil && *user.Name != "" {
			userName = *user.Name
		}

		j.log.Debugf("Sending daily notification to subscribed user ID=%d, phone=%s, name=%s", subscription.UserID, user.PhoneNumber, userName)

		err = j.whatsappService.SendDailyTideNotification(user.PhoneNumber, userName, tidesResponse.Extremes)
		if err != nil {
			j.log.Errorf("Failed to send daily notification to user ID=%d: %v", subscription.UserID, err)
			errorCount++
		} else {
			j.log.Debugf("Successfully sent daily notification to user ID=%d", subscription.UserID)
			successCount++
		}
	}

	j.log.Infof("Daily notifications job completed: %d successful, %d errors out of %d subscribed users", successCount, errorCount, len(subscriptions))

	if errorCount > 0 {
		return successCount, fmt.Errorf("daily notifications job completed with %d errors out of %d subscribed users", errorCount, len(subscriptions))
	}

	return successCount, nil
}
