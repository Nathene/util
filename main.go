package main

import (
	"errors"
	"fmt"
	"github.com/Nathene/util/jwtlib"
	"sync"
	"time"

	"github.com/Nathene/util/ratelimiter"
	"github.com/Nathene/util/bulkhead"
	"github.com/Nathene/util/circuitbreaker"

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

	Bulkhead()

	CircuitBreaker()

	JWT()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Request successful")
	})

	err := e.Start(":8080")
	if err != nil {
		return
	}
}


func Bulkhead() {
	bh := bulkhead.New(3) // Limit to 3 concurrent tasks
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		bh.Execute(func() {
			defer wg.Done()
			fmt.Println("Processing request", time.Now())
			time.Sleep(1 * time.Second)
		})
	}

	wg.Wait()

}

func CircuitBreaker() {
	cb := circuitbreaker.New(3, 5 * time.Second)
		for i := 0; i < 10; i++ {
			err := cb.Call(func() error {
				return errors.New("simulated failure")
			})

			if err != nil {
				fmt.Println("Request failed:", err)
			} else {
				fmt.Println("Request succeeded")
			}
			time.Sleep(1 * time.Second)
		}
}

func JWT() {
	jwtManager := jwtlib.NewJWTManager("secret-salt")

	// Generate token
	token, err := jwtManager.GenerateToken(map[string]interface{}{
		"user": "admin",
	}, time.Hour)

	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}

	fmt.Println("Generated Token:", token)

	// Verify token
	claims, err := jwtManager.VerifyToken(token)
	if err != nil {
		fmt.Println("Invalid token:", err)
		return
	}

	fmt.Println("Token claims:", claims)
}