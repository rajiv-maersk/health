package main

import (
	"health/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// conditionalLogger logs requests only if the status code is not 200
func conditionalLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only log if the status code is not 200
		if c.Writer.Status() != 200 {
			log.Printf("[GIN] %s | %d | %s | %s",
				c.ClientIP(),
				c.Writer.Status(),
				c.Request.Method,
				c.Request.URL.Path,
			)
		}
	}
}

var serviceName string = "accesscontrol"

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Set Gin mode to ReleaseMode for production
	gin.SetMode(gin.ReleaseMode)
	// Initialize Gin router without default middleware
	router := gin.New()

	// Add custom logging middleware
	router.Use(conditionalLogger())

	// Readiness probe handler
	router.GET("/readiness", handlers.ReadinessProbeHandler(serviceName))
	// Liveness probe handler
	router.GET("/liveness", handlers.LivenessProbeHandler(serviceName))

	// Launch the server in a goroutine
	go func() {
		if err := router.Run("localhost:8082"); err != nil {
			// Handle error if server fails to start
			logger.Error("Failed to start server", zap.Error(err))
		}
	}()

	// Keep the main function running
	select {}
}
