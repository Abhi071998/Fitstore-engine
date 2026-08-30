package config

import (
	"log"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload" // Automatically scans and loads the local .env file
)

// Config holds our validated environment settings
type Config struct {
	Env          string
	Port         string
	DatabaseURL  string
	JWTSecret    string
	FrontendURLs []string
}

// Load reads values from the system environment and packages them up
func Load() *Config {
	log.Println("🔍 [CONFIG] Initiating environment variable loading sequence...")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("🚨 [CONFIG] System boot failed: CRITICAL variable 'DATABASE_URL' is missing from the environment.")
	}
	log.Println("✅ [CONFIG] DATABASE_URL parameter successfully mapped from system properties.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("⚠️ [CONFIG] PORT variable not found. Falling back to default: 8080")
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
		log.Println("⚠️ [CONFIG] ENV profile variable empty. Falling back to standard default: development")
	}

	// Ingest the JWT token encryption key signature
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_development_secret_signature_key_123"
		log.Println("⚠️ [CONFIG] JWT_SECRET variable not found. Falling back to an insecure development signature!")
	} else {
		log.Println("✅ [CONFIG] JWT_SECRET signature successfully verified and loaded.")
	}

	// 🌐 Ingest the Frontend Client URL(s) for CORS policy whitelist.
	// Accepts a single origin or a comma-separated list (e.g. local dev + deployed frontend).
	rawFrontendURLs := os.Getenv("FRONTEND_URL")
	if rawFrontendURLs == "" {
		log.Fatal("🚨 [CONFIG] System boot failed: CRITICAL variable 'FRONTEND_URL' is missing from the environment.")
	}
	var frontendURLs []string
	for _, origin := range strings.Split(rawFrontendURLs, ",") {
		origin = strings.TrimSpace(origin)
		origin = strings.TrimSuffix(origin, "/")
		if origin != "" {
			frontendURLs = append(frontendURLs, origin)
		}
	}
	log.Printf("✅ [CONFIG] FRONTEND_URL client origin(s) mapped: %v", frontendURLs)

	log.Printf("🎯 [CONFIG] Configuration binding pipeline complete. Ready to pass structs.")
	return &Config{
		Env:          env,
		Port:         port,
		DatabaseURL:  dbURL,
		JWTSecret:    secret,
		FrontendURLs: frontendURLs,
	}
}
