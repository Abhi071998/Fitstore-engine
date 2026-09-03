package routes

import (
	"log"
	"net/http"

	"tshirt-store/internal/content"
	"tshirt-store/internal/handlers"
	"tshirt-store/internal/middleware" // Import your custom middleware package
	"tshirt-store/internal/orders"

	"github.com/labstack/echo/v4"
)

// SetupRoutes acts as the master map for all API endpoints in the application
func SetupRoutes(
	e *echo.Echo,
	authHandler *handlers.AuthHandler,
	productHandler *handlers.ProductHandler,
	categoryHandler *handlers.CategoryHandler,
	categoryTypeHandler *handlers.CategoryTypeHandler,
	orderHandler *orders.Handler,
	contentHandler *content.Handler,
	jwtSecret []byte, // Injected secret key for token verification
) {
	log.Println("🗺️  [ROUTES] Registering application route mappings...")

	// 🔑 Instantiate the custom authentication middleware
	authRequired := middleware.JWTMiddleware(jwtSecret)

	// 1. Core System Infrastructure Endpoints
	e.GET("/health", func(c echo.Context) error {
		log.Println("🔍 [ROUTES] Health check endpoint touched.")
		return c.JSON(http.StatusOK, map[string]string{
			"status": "online",
			"orm":    "synchronized",
		})
	})

	// 2. Authentication Route Group (Public)
	authGroup := e.Group("/api/auth")
	{
		log.Println("🔒 [ROUTES] Binding authentication subsystem routes...")
		authGroup.POST("/signup", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// 3. Category Routing Blueprint
	categoryGroup := e.Group("/api/categories")
	{
		log.Println("📁 [ROUTES] Binding category configuration routes...")
		// 🔒 Protected: Requires a valid Bearer token signature
		categoryGroup.GET("/getAllCategories", categoryHandler.GetAllCategories, authRequired)

		// 🔒 Protected: Requires a valid Bearer token signature
		categoryGroup.POST("/createCategory", categoryHandler.CreateCategory, authRequired)
		categoryGroup.PUT("/updateCategory/:id", categoryHandler.UpdateCategory, authRequired)
		categoryGroup.DELETE("/deleteCategory/:id", categoryHandler.DeleteCategory, authRequired)
	}

	// 3b. Category Type Routing Blueprint (fixed dropdown list for category creation)
	categoryTypeGroup := e.Group("/api/categoryTypes")
	{
		log.Println("📁 [ROUTES] Binding category type configuration routes...")
		// 🔒 Protected: Requires a valid Bearer token signature
		categoryTypeGroup.GET("/getAllCategoryTypes", categoryTypeHandler.GetAllCategoryTypes, authRequired)

		// 🔒 Protected: Requires a valid Bearer token signature
		categoryTypeGroup.POST("/createCategoryType", categoryTypeHandler.CreateCategoryType, authRequired)
		categoryTypeGroup.PUT("/updateCategoryType/:id", categoryTypeHandler.UpdateCategoryType, authRequired)
		categoryTypeGroup.DELETE("/deleteCategoryType/:id", categoryTypeHandler.DeleteCategoryType, authRequired)
	}

	// 4. Product Routing Blueprint
	productGroup := e.Group("/api/products")
	{
		log.Println("👕 [ROUTES] Binding product core routes...")
		// 🔒 Protected: Requires a valid Bearer token signature
		productGroup.GET("/getAllProducts/:categoryId", productHandler.GetAllProducts, authRequired)

		// 🔒 Protected: Requires a valid Bearer token signature
		productGroup.POST("/createProduct", productHandler.CreateProduct, authRequired)
		productGroup.PUT("/updateProduct/:id", productHandler.UpdateProduct, authRequired)
		productGroup.DELETE("/deleteProduct/:id", productHandler.DeleteProduct, authRequired)
	}

	// 5. Order Routing Blueprint (reads fitstore-core's customer order tables)
	orderGroup := e.Group("/api/orders")
	{
		log.Println("📦 [ROUTES] Binding order review routes...")
		// 🔒 Protected: Requires a valid Bearer token signature
		orderGroup.GET("/pending", orderHandler.GetPendingOrders, authRequired)
		orderGroup.PUT("/:id/approve", orderHandler.ApproveOrder, authRequired)
		orderGroup.PUT("/:id/reject", orderHandler.RejectOrder, authRequired)
		orderGroup.POST("/bulk-approve", orderHandler.BulkApproveOrders, authRequired)
	}

	// 6. Admin Content Routing Blueprint (editable page copy/images, e.g. About Us)
	contentGroup := e.Group("/api/content")
	{
		log.Println("🖼️  [ROUTES] Binding admin content routes...")
		// 🌐 Public: Frontend pages read this to render themselves
		contentGroup.GET("/about-us", contentHandler.GetAboutUs)

		// 🔒 Protected: Requires a valid Bearer token signature
		contentGroup.POST("/about-us", contentHandler.CreateAboutUs, authRequired)
		contentGroup.PUT("/about-us", contentHandler.UpdateAboutUs, authRequired)

		// 🌐 Public: Frontend homepage reads this to render the hero banner
		contentGroup.GET("/hero", contentHandler.GetHero)

		// 🔒 Protected: Requires a valid Bearer token signature
		contentGroup.POST("/hero", contentHandler.CreateHero, authRequired)
		contentGroup.PUT("/hero", contentHandler.UpdateHero, authRequired)
	}

	log.Println("✅ [ROUTES] All backend network endpoints mapped successfully.")
}
