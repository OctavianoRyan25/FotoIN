package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/middleware"
	"fotoin/backend/internal/usecase"
	"fotoin/backend/pkg/response"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PhotoHandler struct {
	photoUsecase *usecase.PhotoUsecase
	s3           *infra.S3Client
}

func NewPhotoHandler(photoUsecase *usecase.PhotoUsecase, s3 *infra.S3Client) *PhotoHandler {
	return &PhotoHandler{photoUsecase: photoUsecase, s3: s3}
}

func (h *PhotoHandler) GetUploadURL(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	galleryIDStr := c.Param("id")
	galleryID, err := uuid.Parse(galleryIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid gallery ID", nil)
	}

	var body struct {
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
	}
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	res, err := h.photoUsecase.GetUploadURL(c.Request().Context(), photographerID, usecase.GetUploadURLInput{
		GalleryID:   galleryID,
		Filename:    body.Filename,
		ContentType: body.ContentType,
	})
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, "Upload URL generated", res)
}

func (h *PhotoHandler) ConfirmUpload(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	var input usecase.ConfirmUploadInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	photo, err := h.photoUsecase.ConfirmUpload(c.Request().Context(), photographerID, input)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Photo registered successfully", photo)
}

func (h *PhotoHandler) ListPhotos(c echo.Context) error {
	galleryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid gallery ID", nil)
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	photos, total, err := h.photoUsecase.ListPhotos(c.Request().Context(), galleryID, page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to fetch photos")
	}

	return response.Paginated(c, "Photos fetched", photos, page, limit, total)
}

func (h *PhotoHandler) DeletePhoto(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	photoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid photo ID", nil)
	}

	if err := h.photoUsecase.Delete(c.Request().Context(), photoID, photographerID); err != nil {
		if err.Error() == "unauthorized" {
			return response.Forbidden(c, "You don't own this photo")
		}
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, "Photo deleted", nil)
}

// ─── Event Photo Endpoints ────────────────────────────────────────────────────

func (h *PhotoHandler) ListEventPhotos(c echo.Context) error {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid event ID", nil)
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	photos, total, err := h.photoUsecase.ListEventPhotos(c.Request().Context(), eventID, page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to fetch photos")
	}

	return response.Paginated(c, "Photos fetched", photos, page, limit, total)
}

func (h *PhotoHandler) GetEventUploadURL(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid event ID", nil)
	}

	var body struct {
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
	}
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	res, err := h.photoUsecase.GetEventUploadURL(c.Request().Context(), photographerID, usecase.GetEventUploadURLInput{
		EventID:     eventID,
		Filename:    body.Filename,
		ContentType: body.ContentType,
	})
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, "Upload URL generated", res)
}

func (h *PhotoHandler) ConfirmEventUpload(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	var input usecase.ConfirmEventUploadInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	photo, err := h.photoUsecase.ConfirmEventUpload(c.Request().Context(), photographerID, input)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Photo registered successfully", photo)
}

// ─── Proxy Upload (backend → S3, bypasses browser CORS) ──────────────────────
// Use when S3 bucket does not have a CORS policy allowing PUT from the browser.
// The browser sends multipart/form-data to this endpoint; the backend proxies
// the file to S3 using server-side credentials (no CORS needed on the bucket).

const maxUploadSize = 50 << 20 // 50 MB

func (h *PhotoHandler) ProxyUpload(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	galleryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid gallery ID", nil)
	}

	if err := c.Request().ParseMultipartForm(maxUploadSize); err != nil {
		return response.BadRequest(c, "File too large or invalid form", err.Error())
	}

	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return response.BadRequest(c, "Missing file field", err.Error())
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Build S3 key — same pattern as presigned flow
	photoID := uuid.New()
	ext := filepath.Ext(header.Filename)
	s3Key := fmt.Sprintf("galleries/%s/photos/%s%s", galleryID, photoID, ext)

	if err := h.s3.UploadObject(s3Key, contentType, file, header.Size); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "S3 upload failed: " + err.Error(),
		})
	}

	// Confirm in DB
	photo, err := h.photoUsecase.ConfirmUpload(c.Request().Context(), photographerID, usecase.ConfirmUploadInput{
		GalleryID: galleryID.String(),
		S3Key:     s3Key,
		Filename:  header.Filename,
		SizeBytes: header.Size,
	})
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Photo uploaded successfully", photo)
}

func (h *PhotoHandler) ProxyEventUpload(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid event ID", nil)
	}

	if err := c.Request().ParseMultipartForm(maxUploadSize); err != nil {
		return response.BadRequest(c, "File too large or invalid form", err.Error())
	}

	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return response.BadRequest(c, "Missing file field", err.Error())
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" || !strings.HasPrefix(contentType, "image/") {
		contentType = "application/octet-stream"
	}

	photoID := uuid.New()
	ext := filepath.Ext(header.Filename)
	s3Key := fmt.Sprintf("events/%s/photos/%s%s", eventID, photoID, ext)

	if err := h.s3.UploadObject(s3Key, contentType, file, header.Size); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "S3 upload failed: " + err.Error(),
		})
	}

	photo, err := h.photoUsecase.ConfirmEventUpload(c.Request().Context(), photographerID, usecase.ConfirmEventUploadInput{
		EventID:   eventID.String(),
		S3Key:     s3Key,
		Filename:  header.Filename,
		SizeBytes: header.Size,
	})
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Photo uploaded successfully", photo)
}
