package handler

import (
	"net/http"

	"fotoin/backend/internal/middleware"
	"fotoin/backend/internal/usecase"
	"fotoin/backend/pkg/response"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// EditProfile — PUT /me/profile
// Updates the authenticated user's display name.
// Security: reads user ID from JWT (not URL param).
func (h *UserHandler) EditProfile(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	var input usecase.UpdateProfileInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	user, err := h.userUsecase.EditProfile(c.Request().Context(), userID, input)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profil berhasil diperbarui",
		"data":    user,
	})
}

// UpdatePassword — PUT /me/password
// Changes password for users who already have one (email-registered or already set via CreatePassword).
func (h *UserHandler) UpdatePassword(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	var input usecase.UpdatePasswordInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := h.userUsecase.UpdatePassword(c.Request().Context(), userID, input); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password berhasil diperbarui",
	})
}

// CreatePassword — POST /me/password
// Sets a password for users who registered via Google OAuth and have no password yet.
func (h *UserHandler) CreatePassword(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	var input usecase.CreatePasswordInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	if err := h.userUsecase.CreatePassword(c.Request().Context(), userID, input); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password berhasil dibuat",
	})
}

// DeleteAccount — DELETE /me
// Permanently deletes the authenticated user's account.
func (h *UserHandler) DeleteAccount(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	if err := h.userUsecase.DeleteAccount(c.Request().Context(), userID); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Akun berhasil dihapus",
	})
}

// ForgotPassword — POST /auth/forgot-password
// Generates an OTP and sends it via email.
func (h *UserHandler) ForgotPassword(c echo.Context) error {
	var input struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}
	if err := c.Validate(input); err != nil {
		return response.BadRequest(c, "Validation failed", err.Error())
	}

	if err := h.userUsecase.ForgetPassword(c.Request().Context(), input.Email); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, "Kode OTP telah dikirim ke email Anda", nil)
}

// ResetPassword — POST /auth/reset-password
// Verifies OTP and sets a new password.
func (h *UserHandler) ResetPassword(c echo.Context) error {
	var input usecase.ResetPasswordInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}
	if err := c.Validate(input); err != nil {
		return response.BadRequest(c, "Validation failed", err.Error())
	}

	if err := h.userUsecase.ResetPassword(c.Request().Context(), input); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, "Password berhasil diubah, silakan login kembali", nil)
}
