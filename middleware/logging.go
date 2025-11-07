package middleware

import (
	"time"

	"github.com/Eursukkul/fiber-booking-system/utils"
	"github.com/gofiber/fiber/v2"
)

type LoggerMiddleware struct {
}

func NewLoggerMiddleware() *LoggerMiddleware {
	return &LoggerMiddleware{}
}

func (m *LoggerMiddleware) Logger(c *fiber.Ctx) error {
	// Generate or retrieve trace_id from header
	traceID := c.Get("X-Trace-ID")
	if traceID == "" {
		traceID = utils.GenerateTraceID()
	}

	// Store trace_id in locals for access in handlers
	c.Locals("trace_id", traceID)

	// Set trace_id in response header
	c.Set("X-Trace-ID", traceID)

	startTime := time.Now()

	// Log incoming request
	utils.LogInfo(traceID, "Incoming request", map[string]interface{}{
		"method": c.Method(),
		"path":   c.Path(),
		"ip":     c.IP(),
	})

	err := c.Next()

	duration := time.Since(startTime)
	statusCode := c.Response().StatusCode()

	// Log response with structured data
	logFields := map[string]interface{}{
		"method":      c.Method(),
		"path":        c.Path(),
		"status":      statusCode,
		"duration_ms": duration.Milliseconds(),
		"ip":          c.IP(),
		"user_agent":  c.Get("User-Agent"),
	}

	// Log at appropriate level based on status code
	if statusCode >= 500 {
		utils.LogError(traceID, err, "Request completed with server error", logFields)
	} else if statusCode >= 400 {
		utils.LogWarn(traceID, "Request completed with client error", logFields)
	} else {
		utils.LogInfo(traceID, "Request completed successfully", logFields)
	}

	return err
}
