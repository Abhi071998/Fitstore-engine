package content

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

// GET /api/content/about-us (Public)
func (h *Handler) GetAboutUs(c echo.Context) error {
	var about AdminContent
	if err := h.DB.First(&about).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "About Us content has not been set up yet"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch About Us content"})
	}
	return c.JSON(http.StatusOK, about)
}

// POST /api/content/about-us (Protected)
// Creates the About Us content. Only one row is expected to ever exist —
// use PUT afterwards to edit it.
func (h *Handler) CreateAboutUs(c echo.Context) error {
	var existing AdminContent
	if err := h.DB.First(&existing).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "About Us content already exists, use PUT to update it"})
	} else if err != gorm.ErrRecordNotFound {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database lookup failed"})
	}

	dto := new(AboutUsDTO)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload format"})
	}

	about := AdminContent{
		AboutUsImg:         dto.AboutUsImg,
		AboutUsTitle:       dto.AboutUsTitle,
		AboutUsDescription: dto.AboutUsDescription,
		AboutUsTagline1:    dto.AboutUsTagline1,
		AboutUsTagline2:    dto.AboutUsTagline2,
		AboutUsTagline3:    dto.AboutUsTagline3,
		AboutUsTagline4:    dto.AboutUsTagline4,
		AboutUsEstbYear:    dto.AboutUsEstbYear,
		AboutUsVisitUs:     dto.AboutUsVisitUs,
		AboutUsEmail:       dto.AboutUsEmail,
	}

	if err := h.DB.Create(&about).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create About Us content"})
	}

	return c.JSON(http.StatusCreated, about)
}

// PUT /api/content/about-us (Protected)
// Updates the single existing About Us row.
func (h *Handler) UpdateAboutUs(c echo.Context) error {
	var existing AdminContent
	if err := h.DB.First(&existing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "About Us content not found, create it first"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database lookup failed"})
	}

	dto := new(AboutUsDTO)
	if err := c.Bind(dto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload format"})
	}

	updates := map[string]interface{}{
		"about_us_img":         dto.AboutUsImg,
		"about_us_title":       dto.AboutUsTitle,
		"about_us_description": dto.AboutUsDescription,
		"about_us_tagline1":    dto.AboutUsTagline1,
		"about_us_tagline2":    dto.AboutUsTagline2,
		"about_us_tagline3":    dto.AboutUsTagline3,
		"about_us_tagline4":    dto.AboutUsTagline4,
		"about_us_estb_year":   dto.AboutUsEstbYear,
		"about_us_visit_us":    dto.AboutUsVisitUs,
		"about_us_email":       dto.AboutUsEmail,
	}

	if err := h.DB.Model(&existing).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update About Us content"})
	}

	var updated AdminContent
	h.DB.First(&updated, existing.ID)
	return c.JSON(http.StatusOK, updated)
}
