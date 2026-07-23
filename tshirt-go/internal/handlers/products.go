package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tshirt-store/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

// POST /api/products/createProduct (Protected)
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	log.Println("📥 [PRODUCT-HANDLER] Processing detailed apparel creation request...")

	dto := new(models.CreateProductDTO)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload format"})
	}

	// Basic validation
	if dto.Name == "" || dto.ProductCode == "" || dto.SKU == "" || dto.CategoryID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name, ProductCode, SKU, and CategoryID are strictly required"})
	}

	// 1. Calculate discount percentage automatically
	var discountPct int
	if dto.MRP > 0 && dto.SellingPrice < dto.MRP {
		discountPct = int(((dto.MRP - dto.SellingPrice) / dto.MRP) * 100)
	}

	// 2. Convert Images slice to JSON string
	imagesJSON, _ := json.Marshal(dto.Images)

	// 3. Convert Specifications map to JSON string
	specsJSON, _ := json.Marshal(dto.Specifications)

	// 4. Build Product entity
	product := models.Product{
		Name:           dto.Name,
		Brand:          dto.Brand,
		Description:    dto.Description,
		ProductCode:    dto.ProductCode,
		SKU:            dto.SKU,
		MRP:            dto.MRP,
		SellingPrice:   dto.SellingPrice,
		Discount:       discountPct,
		CategoryID:     dto.CategoryID,
		Images:         string(imagesJSON),
		Specifications: string(specsJSON),
	}

	// 5. Build Size variants
	for _, sizeDTO := range dto.Sizes {
		product.Sizes = append(product.Sizes, models.ProductSize{
			Size:  sizeDTO.Size,
			Stock: sizeDTO.Stock,
		})
	}

	// 6. Save cleanly in PostgreSQL (GORM will save Product and ProductSizes in a single transaction)
	if err := h.DB.Create(&product).Error; err != nil {
		log.Printf("🚨 [PRODUCT-HANDLER] DB Save error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create product record"})
	}

	log.Printf("🎉 [PRODUCT-HANDLER] Success! Apparel '%s' (Code: %s) created with ID #%d", product.Name, product.ProductCode, product.ID)
	return c.JSON(http.StatusCreated, product)
}

// GET /api/products/getAllProducts (Public)
func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	var products []models.Product
	// Preload Category and Sizes breakdown
	if err := h.DB.Preload("Category").Preload("Sizes").Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch catalog products"})
	}
	return c.JSON(http.StatusOK, products)
}

// DELETE /api/products/deleteProduct/:id (Protected)
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	idParam := c.Param("id")
	productID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID format"})
	}

	// 1. Check if product exists
	var product models.Product
	if err := h.DB.First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database lookup failed"})
	}

	// 2. Perform soft-delete (GORM handles cascading soft-deletes or flags deleted_at)
	if err := h.DB.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete product"})
	}

	log.Printf("🗑️ [PRODUCT-HANDLER] Product ID #%d deleted successfully", productID)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Product deleted successfully",
	})
}
