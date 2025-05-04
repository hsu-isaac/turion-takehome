package handlers

import (
	"telemetry-api/internal/database"
	"telemetry-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	db *database.Database
}

func NewHandlers(db *database.Database) *Handlers {
	return &Handlers{db: db}
}

func (h *Handlers) GetTelemetry(c *fiber.Ctx) error {
	query := &models.TelemetryQuery{
		Page:     1,   // Set default values
		PageSize: 100, // Set default values
	}

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	if query.StartTime.IsZero() || query.EndTime.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "start_time and end_time are required",
		})
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 100
	}

	// Set default subsystem ID if not provided
	if query.SubsystemID == nil {
		defaultID := uint16(1)
		query.SubsystemID = &defaultID
	}

	records, totalCount, err := h.db.GetTelemetry(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch telemetry data",
		})
	}

	response := models.TelemetryResponse{
		Data: records,
		Metadata: struct {
			TotalCount int              `json:"total_count"`
			PageCount  int              `json:"page_count"`
			HasMore    bool             `json:"has_more"`
			TimeRange  models.TimeRange `json:"time_range"`
		}{
			TotalCount: totalCount,
			PageCount:  (totalCount + query.PageSize - 1) / query.PageSize,
			HasMore:    totalCount > query.Page*query.PageSize,
			TimeRange: models.TimeRange{
				Start: query.StartTime,
				End:   query.EndTime,
			},
		},
	}

	return c.JSON(response)
}

func (h *Handlers) GetCurrentTelemetry(c *fiber.Ctx) error {
	record, err := h.db.GetCurrentTelemetry()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch current telemetry",
		})
	}

	return c.JSON(record)
}

func (h *Handlers) GetAnomalies(c *fiber.Ctx) error {
	query := &models.TelemetryQuery{
		Page:     1,   // Set default values
		PageSize: 100, // Set default values
	}

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	if query.StartTime.IsZero() || query.EndTime.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "start_time and end_time are required",
		})
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 100
	}

	// Set default subsystem ID if not provided
	if query.SubsystemID == nil {
		defaultID := uint16(1)
		query.SubsystemID = &defaultID
	}

	anomalies, totalCount, err := h.db.GetAnomalies(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch anomalies",
		})
	}

	response := models.TelemetryResponse{
		Data: anomalies,
		Metadata: struct {
			TotalCount int              `json:"total_count"`
			PageCount  int              `json:"page_count"`
			HasMore    bool             `json:"has_more"`
			TimeRange  models.TimeRange `json:"time_range"`
		}{
			TotalCount: totalCount,
			PageCount:  (totalCount + query.PageSize - 1) / query.PageSize,
			HasMore:    totalCount > query.Page*query.PageSize,
			TimeRange: models.TimeRange{
				Start: query.StartTime,
				End:   query.EndTime,
			},
		},
	}

	return c.JSON(response)
}

func (h *Handlers) GetAggregates(c *fiber.Ctx) error {
	query := new(models.TelemetryAggregationQuery)

	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	metrics, err := h.db.GetAggregatedTelemetry(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch aggregated telemetry",
		})
	}

	return c.JSON(fiber.Map{
		"data": metrics,
		"time_range": models.TimeRange{
			Start: query.StartTime,
			End:   query.EndTime,
		},
	})
}
