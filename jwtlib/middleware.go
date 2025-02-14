package jwtlib

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (j *JWTHandler) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
		}

		isValid, claims := j.Verify(token)
		if !isValid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		// Store claims in the context for future use
		c.Set("user", claims)
		return next(c)
	}
}

func RequireRole(c echo.Context, role string) bool {
	userClaims := c.Get("user").(map[string]interface{})
	return userClaims["role"] == role
}

