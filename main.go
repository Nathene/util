package main

import (
	"time"

	ratelimiter "github.com/Nathene/util/rate_limiter"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

// Global rate limiter instance
var limiter = ratelimiter.New(5, time.Second)

func rateLimitMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !limiter.Allow() {
			return c.JSON(http.StatusTooManyRequests, map[string]string{"error": "Rate limit exceeded"})
		}
		return next(c)
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(rateLimitMiddleware)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Request successful")
	})

	err := e.Start(":8080")
	if err != nil {
		return
	}
}

