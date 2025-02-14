# JWTLib - Simple JWT Authentication for Go

`jwtlib` is a lightweight, easy-to-use JWT (JSON Web Token) library for Go. It provides a simple API for generating, verifying, and managing JWTs, including features like token expiration, middleware for web frameworks, blacklisting, refresh tokens, and role-based access control.

## 🚀 Features
- ✅ **Simple API** - Easily generate and verify JWTs.
- ✅ **Middleware Support** - Works with web frameworks like Echo.
- ✅ **Token Blacklisting** - Supports logout functionality.
- ✅ **Refresh Tokens** - Enables long-term authentication.
- ✅ **Role-Based Access Control (RBAC)** - Assign user roles to control access.
- ✅ **Multi-Algorithm Support** - Works with HS256, RS256, etc.

---

## 🔒 Security Best Practices
##### - Use strong secret keys - Avoid short or predictable keys.
##### - Set expiration times - Use short-lived access tokens (e.g., 15 minutes).
##### - Use refresh tokens securely - Store them server-side or securely on the client.
##### - Blacklist revoked tokens - Prevent users from reusing old tokens.

---

## 📦 Installation

```sh
go get github.com/Nathene/util/jwtlib
```

---

## 🔧 Usage

### 1. Create a JWT Handler
Instantiate a JWTHandler with a secret key:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Nathene/util/jwtlib"
)

func main() {
	// Create a new JWT handler with a custom secret key
	jwtHandler := jwtlib.New("my-super-secret-key")

	// Generate a token with claims
	claims := map[string]interface{}{
		"user": "admin",
		"role": "admin",
	}
	token, err := jwtHandler.Generate(claims, time.Hour)
	if err != nil {
		log.Fatal("Failed to generate token:", err)
	}
	fmt.Println("Generated Token:", token)

	// Verify the token
	isValid, verifiedClaims := jwtHandler.Verify(token)
	if isValid {
		fmt.Println("Token is valid! Claims:", verifiedClaims)
	} else {
		fmt.Println("Token is invalid")
	}
}
```
---

### 2. Middleware for Echo (API Authentication)
Use the built-in middleware to protect routes.

```go
package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/Nathene/util/jwtlib"
)

func main() {
	e := echo.New()
	jwtHandler := jwtlib.New("supersecret")

	// Middleware for authentication
	e.Use(jwtHandler.Middleware)

	e.GET("/protected", func(c echo.Context) error {
		user := c.Get("user").(map[string]interface{})
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome, " + user["user"].(string),
		})
	})

	e.Start(":8080")
}
```

---

### 3. Token Blacklisting (Logout)
Invalidate a token by adding it to a Redis blacklist.

```go
package main

import (
	"fmt"
	"time"

	"github.com/Nathene/util/jwtlib"
)

func main() {
	token := "eyJhbGciOiJIUzI1NiIsIn..." // Example token

	// Blacklist the token for 1 hour
	err := jwtlib.BlacklistToken(token, time.Hour)
	if err != nil {
		fmt.Println("Failed to blacklist token:", err)
		return
	}

	// Check if token is blacklisted
	if jwtlib.IsTokenBlacklisted(token) {
		fmt.Println("Token is blacklisted")
	} else {
		fmt.Println("Token is still valid")
	}
}

```

---

### 4. Refresh Tokens for Persistent Login
Use refresh tokens to generate new access tokens.

```go
package main

import (
	"fmt"
	"time"

	"github.com/Nathene/util/jwtlib"
)

func main() {
	jwtHandler := jwtlib.New("supersecret")

	// Generate a refresh token
	refreshToken, err := jwtHandler.GenerateRefreshToken("user123", 7*24*time.Hour)
	if err != nil {
		fmt.Println("Error generating refresh token:", err)
		return
	}
	fmt.Println("Refresh Token:", refreshToken)
}

```

### 5. Role-Based Access Control (RBAC)
Restrict API access based on user roles.

```go
package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/Nathene/util/jwtlib"
)

func RequireRole(c echo.Context, role string) bool {
	userClaims := c.Get("user").(map[string]interface{})
	return userClaims["role"] == role
}

func main() {
	e := echo.New()
	jwtHandler := jwtlib.New("supersecret")

	e.Use(jwtHandler.Middleware)

	e.GET("/admin", func(c echo.Context) error {
		if !RequireRole(c, "admin") {
			return echo.ErrUnauthorized
		}
		return c.String(200, "Welcome Admin")
	})

	e.Start(":8080")
}

```