package handlers

import (
	"log"
	"net/http"
	"strconv"
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

// PUT /api/categories/updateCategory/:id (Protected)
func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	// Parse ID from URL path parameter
	idParam := c.Param("id")
	categoryID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID format"})
	}

	// 1. Find the target category
	var category models.Category
	if err := h.DB.First(&category, categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database lookup failed"})
	}

	// 2. Bind payload
	dto := new(models.CreateCategoryDTO)
	if err := c.Bind(dto); err != nil || dto.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category name is required for update"})
	}

	// 3. Prevent renaming to another category's existing name
	var existing models.Category
	if err := h.DB.Where("name = ? AND id != ?", dto.Name, categoryID).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Another category with this name already exists"})
	}

	// 4. Update and Save
	category.Name = dto.Name
	if err := h.DB.Save(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update category"})
	}

	log.Printf("✏️ [CATEGORY-HANDLER] Category ID #%d updated to '%s'", category.ID, category.Name)
	return c.JSON(http.StatusOK, category)
}

// DELETE /api/categories/deleteCategory/:id (Protected)
func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	idParam := c.Param("id")
	categoryID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category ID format"})
	}

	// 1. Check if category exists
	var category models.Category
	if err := h.DB.First(&category, categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database lookup failed"})
	}

	// 2. Soft-delete or hard-delete via GORM
	if err := h.DB.Delete(&category).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete category"})
	}

	log.Printf("🗑️ [CATEGORY-HANDLER] Category ID #%d deleted successfully", categoryID)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category deleted successfully",
	})
}
