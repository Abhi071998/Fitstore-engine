package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// CustomClaims reflects the structure encoded during login
type CustomClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTMiddleware returns an Echo middleware function that verifies access tokens
func JWTMiddleware(jwtSecret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Extract the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Println("🔒 [AUTH-MIDDLEWARE] Request rejected: Missing Authorization header")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization token",
				})
			}

			// 2. Validate header format (Bearer <token>)
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				log.Println("🔒 [AUTH-MIDDLEWARE] Request rejected: Malformed Authorization header")
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header must be in format: Bearer <token>",
				})
			}

			tokenString := parts[1]

			// 3. Parse and cryptographically verify the token signature
			claims := &CustomClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Verify HMAC signing algorithm
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
				}
				return jwtSecret, nil
			})

			// 4. Handle invalid signature or expired token
			if err != nil || !token.Valid {
				log.Printf("🔒 [AUTH-MIDDLEWARE] Validation failed: %v", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid or expired token",
				})
			}

			// 5. Save claims into Echo context so handlers can access UserID & Email
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)

			log.Printf("✅ [AUTH-MIDDLEWARE] Verified token for User #%d (%s)", claims.UserID, claims.Email)

			// 6. Token valid! Continue to the target handler
			return next(c)
		}
	}
}
