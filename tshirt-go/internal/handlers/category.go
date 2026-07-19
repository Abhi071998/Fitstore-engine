package handlers

import (
	"log"
	"net/http"
	"tshirt-store/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	DB *gorm.DB
}

// POST /api/categories
func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	log.Println("📥 [CATEGORY-HANDLER] Processing incoming category insertion request...")

	dto := new(models.CreateCategoryDTO)
	if err := c.Bind(dto); err != nil {
		log.Printf("❌ [CATEGORY-HANDLER] Payload structural binding fault: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload format"})
	}

	if dto.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category name is strictly required"})
	}

	// Double-check if this classification already exists
	var existing models.Category
	if err := h.DB.Where("name = ?", dto.Name).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "This category tier already exists"})
	}

	newCategory := models.Category{
		Name: dto.Name,
	}

	// Insert into the database
	log.Println("💾 [CATEGORY-HANDLER] Executing insertion down to PostgreSQL...")
	if err := h.DB.Create(&newCategory).Error; err != nil {
		log.Printf("🚨 [CATEGORY-HANDLER] DB storage routine failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database entry processing fault"})
	}

	log.Printf("🎉 [CATEGORY-HANDLER] Success! Category '%s' inserted with ID #%d", newCategory.Name, newCategory.ID)
	return c.JSON(http.StatusCreated, newCategory)
}

// GET /api/categories
func (h *CategoryHandler) GetAllCategories(c echo.Context) error {
	var categories []models.Category
	if err := h.DB.Preload("Products").Find(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch taxonomy database"})
	}
	return c.JSON(http.StatusOK, categories)
}
