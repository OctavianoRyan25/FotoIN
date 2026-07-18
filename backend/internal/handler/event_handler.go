package handler

import (
	"fotoin/backend/internal/middleware"
	"fotoin/backend/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	eventUC usecase.EventUsecase
}

func NewEventHandler(eventUC usecase.EventUsecase) *EventHandler {
	return &EventHandler{eventUC: eventUC}
}

type CreateEventRequest struct {
	Title       string  `json:"title" validate:"required"`
	Location    string  `json:"location"`
	Description string  `json:"description"`
	EventDate   *string `json:"event_date"`
}

type UpdateEventRequest struct {
	Title       string  `json:"title"`
	Location    string  `json:"location"`
	Description string  `json:"description"`
	EventDate   *string `json:"event_date"`
	IsActive    *bool   `json:"is_active"`
}

func (h *EventHandler) CreateEvent(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	var req CreateEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
		})
	}

	var eventDate *time.Time
	if req.EventDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EventDate)
		if err == nil {
			eventDate = &parsed
		}
	}

	event, err := h.eventUC.CreateEvent(c.Request().Context(), userID, req.Title, req.Location, req.Description, eventDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Event created successfully",
		"data":    event,
	})
}

func (h *EventHandler) GetEvent(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid event ID",
		})
	}

	event, err := h.eventUC.GetEvent(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "Event not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    event,
	})
}

func (h *EventHandler) UpdateEvent(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid event ID",
		})
	}

	var req UpdateEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
		})
	}

	var eventDate *time.Time
	if req.EventDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EventDate)
		if err == nil {
			eventDate = &parsed
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	if err := h.eventUC.UpdateEvent(c.Request().Context(), id, userID, req.Title, req.Location, req.Description, eventDate, isActive); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Event updated successfully",
	})
}

func (h *EventHandler) DeleteEvent(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid event ID",
		})
	}

	if err := h.eventUC.DeleteEvent(c.Request().Context(), id, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Event deleted successfully",
	})
}

func (h *EventHandler) ListMyEvents(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized",
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	events, total, err := h.eventUC.ListEventsByPhotographer(c.Request().Context(), userID, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to fetch events",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    events,
		"meta": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *EventHandler) GetEventBySlug(c echo.Context) error {
	slug := c.Param("slug")

	event, err := h.eventUC.GetEventBySlug(c.Request().Context(), slug)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "Event not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    event,
	})
}

func (h *EventHandler) ListPublicEvents(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")

	events, total, err := h.eventUC.ListPublicEvents(c.Request().Context(), page, limit, search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to fetch events",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    events,
		"meta": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
