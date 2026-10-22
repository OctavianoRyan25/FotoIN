package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"fotoin/backend/internal/domain"
)

type JWTClaims struct {
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "Authorization header missing or invalid",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("APP_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "Invalid token claims",
			})
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		return next(c)
	}
}

func RequireRole(roles ...domain.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, _ := c.Get("user_role").(domain.UserRole)
			for _, r := range roles {
				if role == r {
					return next(c)
				}
			}
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"success": false,
				"message": "You don't have permission to access this resource",
			})
		}
	}
}

func GetUserID(c echo.Context) (uuid.UUID, error) {
	userIDStr, _ := c.Get("user_id").(string)
	return uuid.Parse(userIDStr)
}

func GetUserRole(c echo.Context) (domain.UserRole, error) {
	role, ok := c.Get("user_role").(domain.UserRole)
	if !ok {
		return "", echo.ErrUnauthorized
	}
	return role, nil
}
