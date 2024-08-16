package main

import (
	"net/http"
	"os"
	""errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HealthResponse represents the structure of the health check response
type HealthResponse struct {
	Status       string            `json:"status"`
	ServiceName  string            `json:"service_name"`
	Dependencies map[string]string `json:"dependencies"`
}

// CheckDatabaseStatus verifies if the database is accessible
func CheckDatabaseStatus() error {
	// Simulate a database check failure
	return errors.New("database check failed")
}

func healthCheck(c *gin.Context) {
	// Perform the health checks
	response := HealthResponse{
		Status:       "healthy",
		ServiceName:  "accesscontrol",
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

func main() {
	// Initialize logger
	logger := InitLogger()
	defer logger.Sync() // Flush any buffered log entries

	// Create a new Gin router
	r := gin.Default()

	// Define the health check endpoint
	r.GET("/health", healthCheck)

	// Start the server
	if err := r.Run(":8082"); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

// InitLogger initializes a zap logger with a human-readable format
func InitLogger() *zap.Logger {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(config)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
