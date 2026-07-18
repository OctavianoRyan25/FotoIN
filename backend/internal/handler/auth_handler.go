package handler

import (
	"net/http"
	"os"

	"fotoin/backend/internal/usecase"
	"fotoin/backend/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var input usecase.RegisterInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}
	if err := c.Validate(input); err != nil {
		return response.BadRequest(c, "Validation failed", err.Error())
	}

	res, err := h.authUsecase.Register(c.Request().Context(), input)
	if err != nil {
		if err.Error() == "email already registered" {
			return response.Conflict(c, err.Error())
		}
		return response.InternalError(c, "Registration failed")
	}

	return response.Created(c, "Registration successful", res)
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var input usecase.LoginInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}
	if err := c.Validate(input); err != nil {
		return response.BadRequest(c, "Validation failed", err.Error())
	}

	res, err := h.authUsecase.Login(c.Request().Context(), input)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.OK(c, "Login successful", res)
}

// GET /api/v1/auth/google
// Redirects the browser to the Google consent screen.
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	// Generate a random CSRF state
	state := uuid.New().String()
	url := h.authUsecase.GoogleLoginURL(state)

	// Redirect the browser directly to Google
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GET /api/v1/auth/google/callback
// Google redirects here with ?code=...&state=...
// We exchange the code for user info, then redirect the client app with the JWT.
func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	errParam := c.QueryParam("error")

	if errParam != "" {
		return c.Redirect(http.StatusTemporaryRedirect,
			frontendURL()+"?error=google_auth_denied")
	}
	if code == "" || state == "" {
		return response.BadRequest(c, "Missing code or state", nil)
	}

	res, err := h.authUsecase.GoogleCallback(c.Request().Context(), code, state)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect,
			frontendURL()+"?error="+err.Error())
	}

	// Redirect client app — embed token in fragment so it's not sent to server
	redirectURL := frontendURL() + "/auth/callback#token=" + res.AccessToken
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// frontendURL returns the client app base URL from env, defaulting to localhost.
func frontendURL() string {
	url := os.Getenv("FRONTEND_URL")
	if url == "" {
		return "http://localhost:3000"
	}
	return url
}
