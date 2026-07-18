package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectORM(databaseURL string) *gorm.DB {
	log.Println("🔌 [DATABASE] Attempting handshake initialization with PostgreSQL instance...")

	dialector := postgres.Open(databaseURL)

	db, err := gorm.Open(dialector, &gorm.Config{
		// Info level prints *every single raw SQL query* your backend fires straight into your terminal
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("🚨 [DATABASE] Handshake connection failed completely. Error trace: %v", err)
	}
	log.Println("✅ [DATABASE] Handshake validation passed. GORM dialer established.")

	// Extract generic sql.DB object to configure pool properties
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("🚨 [DATABASE] Failed to extract base SQL lifecycle controller. Error trace: %v", err)
	}

	log.Println("⚙️ [DATABASE] Restructuring connection pool bounds...")
	sqlDB.SetMaxIdleConns(10)
	log.Println("📥 [DATABASE] Max Idle Connections configured to: 10")
	sqlDB.SetMaxOpenConns(100)
	log.Println("📥 [DATABASE] Max Open Connections configured to: 100")
	sqlDB.SetConnMaxLifetime(time.Hour)
	log.Printf("📥 [DATABASE] Max Connection Lifecycle Duration locked: %v", time.Hour)

	log.Println("🐘 [DATABASE] PostgreSQL connection manager successfully bound to runtime context.")
	return db
}
