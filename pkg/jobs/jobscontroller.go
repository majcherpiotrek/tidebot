package jobs

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type JobsController struct {
	jobsService JobsService
	apiKey      string
	log         echo.Logger
}

func NewJobsController(jobsService JobsService, apiKey string, log echo.Logger) *JobsController {
	return &JobsController{
		jobsService: jobsService,
		apiKey:      apiKey,
		log:         log,
	}
}

func (jc *JobsController) RegisterRoutes(e *echo.Echo) {
	jobsGroup := e.Group("/jobs")
	jobsGroup.Use(jc.apiKeyMiddleware)

	jobsGroup.POST("/send-tide-extremes", jc.SendTideExtremesToAllUsers)
	jobsGroup.POST("/v2/send-daily-notifications", jc.SendDailyNotifications)
}

func (jc *JobsController) apiKeyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get("X-API-Key")
		if apiKey == "" {
			// Also check Authorization header with Bearer format
			auth := c.Request().Header.Get("Authorization")
			if after, ok := strings.CutPrefix(auth, "Bearer "); ok {
				apiKey = after
			}
		}

		if apiKey != jc.apiKey {
			jc.log.Warnf("Invalid API key attempt from %s", c.RealIP())
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid API key",
			})
		}

		return next(c)
	}
}

func (jc *JobsController) SendTideExtremesToAllUsers(c echo.Context) error {
	jc.log.Info("Received request to send tide extremes to all users")

	err := jc.jobsService.SendTideExtremesToAllUsers()
	if err != nil {
		jc.log.Errorf("Failed to send tide extremes to all users: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to send tide extremes to all users",
			"error":   err.Error(),
		})
	}

	jc.log.Info("Successfully completed sending tide extremes to all users")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Tide extremes sent to all users successfully",
	})
}

func (jc *JobsController) SendDailyNotifications(c echo.Context) error {
	jc.log.Info("Received request to send daily tide notifications (v2)")

	successCount, err := jc.jobsService.SendDailyNotificationsV2()
	if err != nil {
		jc.log.Errorf("Failed to send daily notifications: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to send daily notifications",
			"error":   err.Error(),
		})
	}

	jc.log.Info("Successfully completed sending daily notifications")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":       "success",
		"successCount": strconv.Itoa(successCount),
		"message":      "Daily notifications sent successfully",
	})
}

