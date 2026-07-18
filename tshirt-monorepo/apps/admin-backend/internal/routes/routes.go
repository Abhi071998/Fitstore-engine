package routes

import (
	"log"
	"net/http"

	"tshirt-store/internal/handlers"

	"github.com/labstack/echo/v4"
)

// SetupRoutes acts as the master map for all API endpoints in the application
func SetupRoutes(e *echo.Echo, authHandler *handlers.AuthHandler) {
	log.Println("🗺️  [ROUTES] Registering application route mappings...")

	// 1. Core System Infrastructure Endpoints
	e.GET("/health", func(c echo.Context) error {
		log.Println("🔍 [ROUTES] Health check endpoint touched.")
		return c.JSON(http.StatusOK, map[string]string{
			"status": "online",
			"orm":    "synchronized",
		})
	})

	// 2. Authentication Route Group
	authGroup := e.Group("/api/auth")
	{
		log.Println("🔒 [ROUTES] Binding authentication subsystem routes...")
		authGroup.POST("/signup", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// You can easily add product groups here later without ever touching main.go again:
	// productGroup := e.Group("/api/products")
	// ...

	log.Println("✅ [ROUTES] All backend network endpoints mapped successfully.")
}
