package handlers

import (
	"log"
	"net/http"
	"strconv"
	"tshirt-store/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CategoryTypeHandler struct {
	DB *gorm.DB
}

// POST /api/categoryTypes/createCategoryType
func (h *CategoryTypeHandler) CreateCategoryType(c echo.Context) error {
	log.Println("📥 [CATEGORY-TYPE-HANDLER] Processing incoming category type insertion request...")

	dto := new(models.CreateCategoryTypeDTO)
	if err := c.Bind(dto); err != nil {
		log.Printf("❌ [CATEGORY-TYPE-HANDLER] Payload structural binding fault: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload format"})
	}

	if dto.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category type name is strictly required"})
	}

	// Double-check if this classification already exists
	var existing models.CategoryType
	if err := h.DB.Where("name = ?", dto.Name).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "This category type already exists"})
	}

	newCategoryType := models.CategoryType{
		Name: dto.Name,
	}

	log.Println("💾 [CATEGORY-TYPE-HANDLER] Executing insertion down to PostgreSQL...")
	if err := h.DB.Create(&newCategoryType).Error; err != nil {
		log.Printf("🚨 [CATEGORY-TYPE-HANDLER] DB storage routine failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database entry processing fault"})
	}

	log.Printf("🎉 [CATEGORY-TYPE-HANDLER] Success! Category type '%s' inserted with ID #%d", newCategoryType.Name, newCategoryType.ID)
	return c.JSON(http.StatusCreated, newCategoryType)
}

// GET /api/categoryTypes/getAllCategoryTypes
func (h *CategoryTypeHandler) GetAllCategoryTypes(c echo.Context) error {
	var categoryTypes []models.CategoryType
	if err := h.DB.Find(&categoryTypes).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch category type database"})
	}
	return c.JSON(http.StatusOK, categoryTypes)
}

// PUT /api/categoryTypes/updateCategoryType/:id (Protected)
func (h *CategoryTypeHandler) UpdateCategoryType(c echo.Context) error {
	idParam := c.Param("id")
	categoryTypeID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category type ID format"})
	}

	var categoryType models.CategoryType
	if err := h.DB.First(&categoryType, categoryTypeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Category type not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database lookup failed"})
	}

	dto := new(models.CreateCategoryTypeDTO)
	if err := c.Bind(dto); err != nil || dto.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category type name is required for update"})
	}

	var existing models.CategoryType
	if err := h.DB.Where("name = ? AND id != ?", dto.Name, categoryTypeID).First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Another category type with this name already exists"})
	}

	categoryType.Name = dto.Name

	// Rename the type and cascade the new name to every linked Category atomically
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&categoryType).Error; err != nil {
			return err
		}
		return tx.Model(&models.Category{}).
			Where("category_type_id = ?", categoryTypeID).
			Update("name", categoryType.Name).Error
	})
	if err != nil {
		log.Printf("🚨 [CATEGORY-TYPE-HANDLER] Cascade update failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update category type"})
	}

	log.Printf("✏️ [CATEGORY-TYPE-HANDLER] Category type ID #%d updated to '%s' (cascaded to linked categories)", categoryType.ID, categoryType.Name)
	return c.JSON(http.StatusOK, categoryType)
}

// DELETE /api/categoryTypes/deleteCategoryType/:id (Protected)
func (h *CategoryTypeHandler) DeleteCategoryType(c echo.Context) error {
	idParam := c.Param("id")
	categoryTypeID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category type ID format"})
	}

	var categoryType models.CategoryType
	if err := h.DB.First(&categoryType, categoryTypeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Category type not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database lookup failed"})
	}

	if err := h.DB.Delete(&categoryType).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete category type"})
	}

	log.Printf("🗑️ [CATEGORY-TYPE-HANDLER] Category type ID #%d deleted successfully", categoryTypeID)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category type deleted successfully",
	})
}
