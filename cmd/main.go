package main

import (
	"log"
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

	 //_ "github.com/Eursukkul/fiber-booking-system/docs"
)
// @title Fiber Booking System API
// @version 1.0
// @description Booking system API with Fiber framework
// @host localhost:3000
// @BasePath /api
func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := fiber.New()

	//Allow all origins
	app.Use(cors.New(cors.Config{
        AllowOrigins: "*",                // Allow all origins
        AllowHeaders: "Origin, Content-Type, Accept",
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

	var wg sync.WaitGroup
	bookingUsecase.BackgroundTaskBooking(&wg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		wg.Wait()

		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error shutting down server: %v", err)
		}

		log.Println("Server shut down gracefully")
	}()

	log.Printf("Server is running on %s", config.Port)
	if err := app.Listen(config.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
