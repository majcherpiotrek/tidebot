package jobs

import (
	"fmt"
	"tidebot/pkg/users/services"
	"tidebot/pkg/whatsapp"
	"tidebot/pkg/worldtides"
	"time"

	"github.com/labstack/echo/v4"
)

type JobsService interface {
	SendTideExtremesToAllUsers() error
}

type jobsServiceImpl struct {
	userService      services.UserService
	whatsappService  whatsapp.WhatsAppService
	worldTidesClient worldtides.WorldTidesClient
	log              echo.Logger
}

func NewJobsService(
	userService services.UserService,
	whatsappService whatsapp.WhatsAppService,
	worldTidesClient worldtides.WorldTidesClient,
	log echo.Logger,
) JobsService {
	return &jobsServiceImpl{
		userService:      userService,
		whatsappService:  whatsappService,
		worldTidesClient: worldTidesClient,
		log:              log,
	}
}

func (j *jobsServiceImpl) SendTideExtremesToAllUsers() error {
	j.log.Info("Starting job: Send tide extremes to all users")

	// Get today's date
	today := time.Now().Format("2006-01-02")
	j.log.Debugf("Fetching tide extremes for date: %s", today)

	// Fetch tide extremes from WorldTides API
	tidesResponse, err := j.worldTidesClient.GetTidalExtremesForDay(today)
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

	// Get all registered users
	users, err := j.userService.GetAllUsers()
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	j.log.Infof("Found %d registered users to send tide extremes", len(users))

	// Send tide extremes to each user
	successCount := 0
	errorCount := 0

	for _, user := range users {
		j.log.Debugf("Sending tide extremes to user ID=%d, phone=%s", user.ID, user.PhoneNumber)

		err := j.whatsappService.SendTideExtremesMessage(user.PhoneNumber, tidesResponse.Extremes, today)
		if err != nil {
			j.log.Errorf("Failed to send tide extremes to user ID=%d: %v", user.ID, err)
			errorCount++
		} else {
			j.log.Debugf("Successfully sent tide extremes to user ID=%d", user.ID)
			successCount++
		}
	}

	j.log.Infof("Job completed: %d successful, %d errors out of %d users", successCount, errorCount, len(users))

	if errorCount > 0 {
		return fmt.Errorf("job completed with %d errors out of %d users", errorCount, len(users))
	}

	return nil
}

