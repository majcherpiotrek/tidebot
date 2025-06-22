package jobs

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type JobsController struct {
	jobsService JobsService
	log         echo.Logger
}

func NewJobsController(jobsService JobsService, log echo.Logger) *JobsController {
	return &JobsController{
		jobsService: jobsService,
		log:         log,
	}
}

func (jc *JobsController) RegisterRoutes(e *echo.Echo) {
	jobsGroup := e.Group("/jobs")
	
	jobsGroup.POST("/send-tide-extremes", jc.SendTideExtremesToAllUsers)
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