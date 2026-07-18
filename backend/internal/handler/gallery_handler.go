package handler

import (
	"fotoin/backend/internal/middleware"
	"fotoin/backend/internal/usecase"
	"fotoin/backend/pkg/response"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type GalleryHandler struct {
	galleryUsecase *usecase.GalleryUsecase
}

func NewGalleryHandler(galleryUsecase *usecase.GalleryUsecase) *GalleryHandler {
	return &GalleryHandler{galleryUsecase: galleryUsecase}
}

func (h *GalleryHandler) Create(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	var body struct {
		Title         string  `json:"title"`
		EventDate     string  `json:"event_date"` // ISO8601 string from frontend
		Location      string  `json:"location"`
		Description   string  `json:"description"`
		PricePerPhoto float64 `json:"price_per_photo"`
	}

	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	eventDate, err := time.Parse("2006-01-02", body.EventDate)
	if err != nil {
		return response.BadRequest(c, "Invalid event_date format (use YYYY-MM-DD)", nil)
	}

	gallery, err := h.galleryUsecase.Create(c.Request().Context(), photographerID, usecase.CreateGalleryInput{
		Title:         body.Title,
		EventDate:     eventDate,
		Location:      body.Location,
		Description:   body.Description,
		PricePerPhoto: body.PricePerPhoto,
	})
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Gallery created successfully", gallery)
}

func (h *GalleryHandler) List(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	galleries, total, err := h.galleryUsecase.ListByPhotographer(c.Request().Context(), photographerID, page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to fetch galleries")
	}

	return response.Paginated(c, "Galleries fetched", galleries, page, limit, total)
}

func (h *GalleryHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid gallery ID", nil)
	}

	gallery, err := h.galleryUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.NotFound(c, "Gallery not found")
	}

	return response.OK(c, "Gallery fetched", gallery)
}

func (h *GalleryHandler) Update(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid gallery ID", nil)
	}

	var input usecase.UpdateGalleryInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	gallery, err := h.galleryUsecase.Update(c.Request().Context(), id, photographerID, input)
	if err != nil {
		if err.Error() == "unauthorized: you don't own this gallery" {
			return response.Forbidden(c, err.Error())
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, "Gallery updated", gallery)
}

func (h *GalleryHandler) Delete(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid gallery ID", nil)
	}

	if err := h.galleryUsecase.Delete(c.Request().Context(), id, photographerID); err != nil {
		if err.Error() == "unauthorized" {
			return response.Forbidden(c, err.Error())
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, "Gallery deleted", nil)
}

// GetPublicGallery — no auth required, for client app
func (h *GalleryHandler) GetPublicGallery(c echo.Context) error {
	slug := c.Param("slug")

	gallery, err := h.galleryUsecase.GetBySlug(c.Request().Context(), slug, nil)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.OK(c, "Gallery fetched", gallery)
}

func (h *GalleryHandler) ListPublic(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	galleries, total, err := h.galleryUsecase.ListPublic(c.Request().Context(), page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to fetch public galleries")
	}

	return response.Paginated(c, "Public galleries fetched", galleries, page, limit, total)
}
