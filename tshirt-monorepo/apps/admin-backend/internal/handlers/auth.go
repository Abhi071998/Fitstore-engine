package handlers

import (
	"log"
	"net/http"

	"tshirt-store/internal/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB        *gorm.DB
	JWTSecret []byte
}

// POST /api/auth/signup
func (h *AuthHandler) Register(c echo.Context) error {
	log.Println("📥 [AUTH-HANDLER] Received incoming user registration request...")

	// 1. Bind incoming JSON body parameters straight into our DTO
	dto := new(models.RegisterDTO)
	if err := c.Bind(dto); err != nil {
		log.Printf("❌ [AUTH-HANDLER] Registration failed: Payload structural binding fault: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload format"})
	}
	log.Printf("🔍 [AUTH-HANDLER] Payload parsed. Attempting validation for target email: %s", dto.Email)

	// 2. Simple fallback request sanitization check
	if dto.Email == "" || dto.Password == "" || dto.Name == "" {
		log.Println("❌ [AUTH-HANDLER] Registration failed: Missing critical credentials fields inside DTO body.")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "All input fields (email, name, password) are strictly required"})
	}

	// 3. Check if the user email account already exists in the database
	var existingUser models.User
	result := h.DB.Where("email = ?", dto.Email).First(&existingUser)
	if result.Error == nil {
		log.Printf("⚠️ [AUTH-HANDLER] Registration rejected: Account email %s is already taken.", dto.Email)
		return c.JSON(http.StatusConflict, map[string]string{"error": "This email is already registered"})
	}

	// 4. Securely encrypt/hash raw password characters via Bcrypt
	log.Println("🔑 [AUTH-HANDLER] Encrypting password safety signatures via Bcrypt...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 12) // 12 cost factor provides secure production balancing
	if err != nil {
		log.Printf("🚨 [AUTH-HANDLER] Crypto processing failed entirely: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to safely process account credentials"})
	}

	// 5. Build our database entity profile mapping properties
	newUser := models.User{
		Email:    dto.Email,
		Name:     dto.Name,
		Password: string(hashedPassword),
	}

	// 6. Write record cleanly down to Postgres via ORM layer
	log.Println("💾 [AUTH-HANDLER] Saving user record down to PostgreSQL database...")
	if err := h.DB.Create(&newUser).Error; err != nil {
		log.Printf("🚨 [AUTH-HANDLER] DB storage routine failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database entry routine processing fault"})
	}

	log.Printf("🎉 [AUTH-HANDLER] Success! Account for %s created cleanly with User ID #%d.", newUser.Email, newUser.ID)

	// Send clean response back (User struct automatically hides the password field because of the JSON tag `json:"-"`)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Registration successful!",
		"user":    newUser,
	})
}

// POST /api/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	// We'll implement the JWT creation token logic right after validating signup works!
	return c.JSON(http.StatusNotImplemented, map[string]string{"message": "Login logic coming up next"})
}
