package routes

import (
	"log"
	"net/http"

	"tshirt-store/internal/handlers"

	"github.com/labstack/echo/v4"
)

// SetupRoutes acts as the master map for all API endpoints in the application
func SetupRoutes(e *echo.Echo, authHandler *handlers.AuthHandler,
	// productHandler *handlers.ProductHandler,
	categoryHandler *handlers.CategoryHandler) {
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

	// =========================================================================
	// 🆕 ADDED SECTION: Category Routing Blueprint
	// =========================================================================
	categoryGroup := e.Group("/api/categories")
	{
		log.Println("📁 [ROUTES] Binding category configuration routes...")
		categoryGroup.POST("createCategory", categoryHandler.CreateCategory)
		categoryGroup.GET("getAllCategories", categoryHandler.GetAllCategories)
	}

	// =========================================================================
	// 🆕 ADDED SECTION: Product Routing Blueprint
	// =========================================================================
	// productGroup := e.Group("/api/products")
	// {
	// 	log.Println("👕 [ROUTES] Binding product core routes...")
	// 	productGroup.POST("createProduct", productHandler.CreateProduct)
	// 	productGroup.GET("getAllProducts", productHandler.GetAllProducts)
	// }
	log.Println("✅ [ROUTES] All backend network endpoints mapped successfully.")
}
