package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HealthResponse struct {
	Status       string            `json:"status"`
	ServiceName  string            `json:"service_name"`
	Dependencies map[string]string `json:"dependencies"`
}

// CheckDatabaseStatus verifies if the database is accessible
func CheckDatabaseStatus() error {
	return errors.New("Database connection failed")
}

// ReadinessProbeHandler returns a Gin handler function that checks the health of the service
func ReadinessProbeHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := HealthResponse{
			Status:       "healthy",
			ServiceName:  serviceName,
			Dependencies: make(map[string]string),
		}

		// Check Database status
		databaseStatus := "healthy"
		if err := CheckDatabaseStatus(); err != nil {
			databaseStatus = "unhealthy: " + err.Error()
			if response.Status != "unhealthy" {
				response.Status = "unhealthy"
			}
			zap.L().Error("Database status check failed", zap.Error(err))
		}
		response.Dependencies["Database"] = databaseStatus

		// Set the appropriate HTTP status code based on overall health
		if response.Status == "unhealthy" {
			c.JSON(http.StatusInternalServerError, response) // 500 Internal Server Error
		} else {
			c.JSON(http.StatusOK, response) // 200 OK
		}
	}
}

// LivenessProbeHandler returns a Gin handler function that checks if the service is alive
func LivenessProbeHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a response with only the service's basic health status
		response := HealthResponse{
			Status:       "healthy",
			ServiceName:  serviceName,
			Dependencies: make(map[string]string),
		}

		// Respond with 200 OK to indicate the service is alive
		c.JSON(http.StatusOK, response) // 200 OK
	}
}
