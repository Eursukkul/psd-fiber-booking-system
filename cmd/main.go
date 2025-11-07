package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Eursukkul/fiber-booking-system/config"
	"github.com/Eursukkul/fiber-booking-system/handler"
	"github.com/Eursukkul/fiber-booking-system/middleware"
	"github.com/Eursukkul/fiber-booking-system/repository"
	"github.com/Eursukkul/fiber-booking-system/router"
	"github.com/Eursukkul/fiber-booking-system/usecase"
	"github.com/Eursukkul/fiber-booking-system/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
	//_ "github.com/Eursukkul/fiber-booking-system/docs"
)

// @title Fiber Booking System API
// @version 1.0
// @description Booking system API with Fiber framework
// @host localhost:3000
// @BasePath /api
func main() {
	// Initialize structured logger
	logger := utils.InitLogger()
	logger.Info("Starting Fiber Booking System API")

	config, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load config")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			traceID, ok := c.Locals("trace_id").(string)
			if !ok {
				traceID = utils.GenerateTraceID()
			}

			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			utils.LogError(traceID, err, "Error handling request", map[string]interface{}{
				"path":   c.Path(),
				"method": c.Method(),
				"status": code,
			})

			return c.Status(code).JSON(fiber.Map{
				"error":    err.Error(),
				"trace_id": traceID,
			})
		},
	})

	//Allow all origins
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins
		AllowHeaders: "Origin, Content-Type, Accept, X-Trace-ID",
		AllowMethods: "GET,POST,DELETE",
	}))

	loggerMiddleware := middleware.NewLoggerMiddleware()
	// authMiddleware := middleware.NewAuthMiddleware()

	bookingRepo := repository.NewMockBookingRepository()
	cache := utils.NewInMemoryCache()

	bookingUsecase := usecase.NewBookingUsecase(bookingRepo, cache)
	bookingHandler := handler.NewBookingHandler(bookingUsecase)

	router.SetupRoutes(app, bookingHandler, loggerMiddleware)

	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "psd-fiber-booking-system",
		})
	})

	var wg sync.WaitGroup
	bookingUsecase.BackgroundTaskBooking(&wg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("Shutting down server...")

		wg.Wait()

		if err := app.Shutdown(); err != nil {
			logger.WithError(err).Error("Error shutting down server")
		}

		logger.Info("Server shut down gracefully")
	}()

	logger.WithFields(logrus.Fields{
		"port": config.Port,
		"env":  os.Getenv("ENV"),
	}).Info("Server is running")

	if err := app.Listen(config.Port); err != nil {
		logger.WithError(err).Fatal("Error starting server")
	}

}
