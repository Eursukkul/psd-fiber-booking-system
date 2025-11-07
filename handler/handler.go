package handler

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/Eursukkul/fiber-booking-system/dto"
	"github.com/Eursukkul/fiber-booking-system/usecase"
	"github.com/Eursukkul/fiber-booking-system/utils"
	"github.com/gofiber/fiber/v2"
)

type (
	BookingHandler struct {
		BookingUsecase usecase.BookingUsecase
	}
)

func NewBookingHandler(BookingUsecase usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{BookingUsecase: BookingUsecase}
}

// getTraceID retrieves trace_id from context
func getTraceID(c *fiber.Ctx) string {
	traceID, ok := c.Locals("trace_id").(string)
	if !ok {
		return utils.GenerateTraceID()
	}
	return traceID
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new booking with user_id, service_id, and price
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body dto.BookingRequest true "Booking Request"
// @Success 201 {object} dto.SwaggerResponse{data=dto.BookingResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Router /bookings [post]
func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	traceID := getTraceID(c)

	var req dto.BookingRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogError(traceID, err, "Failed to parse request body", map[string]interface{}{
			"path": c.Path(),
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  "Invalid request body",
			"trace_id": traceID,
		})
	}

	utils.LogInfo(traceID, "Creating booking", map[string]interface{}{
		"user_id":    req.UserID,
		"service_id": req.ServiceID,
		"price":      req.Price,
	})

	booking, err := h.BookingUsecase.CreateBooking(req)
	if err != nil {
		utils.LogError(traceID, err, "Failed to create booking", map[string]interface{}{
			"user_id":    req.UserID,
			"service_id": req.ServiceID,
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":  err.Error(),
			"trace_id": traceID,
		})
	}

	if booking.Price > 50000 {
		utils.LogInfo(traceID, "High value booking detected, starting async credit check", map[string]interface{}{
			"booking_id": booking.ID,
			"price":      booking.Price,
		})

		go func(id int, traceID string) {
			//Random status confirm or rejected
			status := "confirmed"
			if rand.Intn(2) == 0 {
				status = "rejected"
			}
			time.Sleep(time.Second * 1) //Delay 1 second
			//Update status
			err := h.BookingUsecase.UpdateBookingStatus(id, status)
			if err != nil {
				utils.LogError(traceID, err, "Failed to update booking status", map[string]interface{}{
					"booking_id": id,
					"status":     status,
				})
			} else {
				utils.LogInfo(traceID, "Booking status updated", map[string]interface{}{
					"booking_id": id,
					"status":     status,
				})
			}
		}(booking.ID, traceID)
	}

	utils.LogInfo(traceID, "Booking created successfully", map[string]interface{}{
		"booking_id": booking.ID,
	})

	return c.Status(fiber.StatusCreated).JSON(booking)
}

// GetBookingByID godoc
// @Summary Get a booking by ID
// @Description Get booking details by ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} dto.SwaggerResponse{data=dto.BookingResponse}
// @Failure 404 {object} dto.ErrorResponse
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBookingByID(c *fiber.Ctx) error {
	traceID := getTraceID(c)

	id, err := c.ParamsInt("id")
	if err != nil {
		utils.LogError(traceID, err, "Invalid booking ID parameter", nil)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  "Invalid ID",
			"trace_id": traceID,
		})
	}

	booking, err := h.BookingUsecase.GetBookingByID(id)
	if err != nil {
		utils.LogWarn(traceID, "Booking not found", map[string]interface{}{
			"booking_id": id,
		})
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message":  "Booking not found",
			"trace_id": traceID,
		})
	}

	utils.LogInfo(traceID, "Booking retrieved successfully", map[string]interface{}{
		"booking_id": id,
	})

	return c.Status(fiber.StatusOK).JSON(booking)
}

// GetAllBookings godoc
// @Summary Get all bookings
// @Description Get all bookings with optional sorting and filtering
// @Tags bookings
// @Accept json
// @Produce json
// @Param sort query string false "Sort by (price or date)"
// @Param high-value query bool false "Filter high value bookings"
// @Success 200 {object} dto.SwaggerResponse{data=[]dto.BookingResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Router /bookings [get]
func (h *BookingHandler) GetAllBookings(c *fiber.Ctx) error {
	traceID := getTraceID(c)

	sortBy := c.Query("sort", "id")
	highValue := c.Query("high-value")

	utils.LogDebug(traceID, "Fetching all bookings", map[string]interface{}{
		"sort_by":    sortBy,
		"high_value": highValue,
	})

	bookings, err := h.BookingUsecase.GetAllBookings(sortBy, highValue)
	if err != nil {
		utils.LogError(traceID, err, "Failed to retrieve bookings", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to retrieve bookings",
			"trace_id": traceID,
		})
	}

	utils.LogInfo(traceID, "Retrieved all bookings", map[string]interface{}{
		"count": len(bookings),
	})

	return c.JSON(bookings)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel an existing booking by ID
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} dto.SwaggerResponse
// @Failure 400,404 {object} dto.ErrorResponse
// @Router /bookings/{id} [delete]
func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	traceID := getTraceID(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		utils.LogError(traceID, err, "Invalid booking ID", nil)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  "Invalid booking ID",
			"trace_id": traceID,
		})
	}

	// ตรวจสอบสถานะก่อนยกเลิก
	booking, err := h.BookingUsecase.GetBookingByID(id)
	if err != nil {
		utils.LogWarn(traceID, "Booking not found for cancellation", map[string]interface{}{
			"booking_id": id,
		})
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message":  "Booking not found",
			"trace_id": traceID,
		})
	}

	if booking.Status == "confirmed" {
		utils.LogWarn(traceID, "Attempted to cancel confirmed booking", map[string]interface{}{
			"booking_id": id,
			"status":     booking.Status,
		})
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":  "Cannot cancel a confirmed booking",
			"trace_id": traceID,
		})
	}

	// ยกเลิกการจอง
	err = h.BookingUsecase.CancelBooking(id)
	if err != nil {
		utils.LogError(traceID, err, "Failed to cancel booking", map[string]interface{}{
			"booking_id": id,
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message":  err.Error(),
			"trace_id": traceID,
		})
	}

	utils.LogInfo(traceID, "Booking canceled successfully", map[string]interface{}{
		"booking_id": id,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Booking canceled successfully",
		"trace_id": traceID,
	})
}
