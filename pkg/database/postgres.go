package database

import (
	"fmt"
	"log"
	"time"

	"github.com/impk123/Inventory-Management-Mini-System/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitPostgres(dsn string) (*gorm.DB, error) {
	var err error

	// Config Logger
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to Postgres
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	// Config connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return DB, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB) error {

	if err := db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Stock{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")

	// Create indexes if they don't exist
	createIndexes(db)

	return nil
}

// CreateIndexes creates indexes for performance
func createIndexes(db *gorm.DB) {
	// User indexes
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
		CREATE INDEX IF NOT EXISTS idx_users_last_login ON users(last_login_at);
	`)

	// Product indexes
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
		CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
		CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);
		CREATE INDEX IF NOT EXISTS idx_products_quantity ON products(quantity);
		CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);
	`)

	// Stock  indexes
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_stock_product_id ON stock(product_id);
		CREATE INDEX IF NOT EXISTS idx_stock_created_at ON stock(created_at);
		CREATE INDEX IF NOT EXISTS idx_stock_type ON stock(type);
		CREATE INDEX IF NOT EXISTS idx_stock_reference ON stock(reference);
	`)

	log.Println("✅ Database indexes created/verified")
}

// HealthCheck checks database connection
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Check if users table exists
	var result int
	err = db.Raw("SELECT 1 FROM users LIMIT 1").Scan(&result).Error
	if err != nil {
		return fmt.Errorf("database table check failed: %w", err)
	}

	return nil
}

// CloseDB closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// SeedData seeds initial data for development
func SeedData(db *gorm.DB) error {
	// Check if users already exist
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
		// Create admin user
		admin := models.User{
			Email:       "admin@inventory.com",
			Password:    "$2a$12$K9S6N8Uv7iWQYH4lX2mZFe.8cJg3qR5sT", // password: admin123
			Name:        "Admin User",
			Role:        "admin",
			IsActive:    true,
			LastLoginAt: time.Now(),
		}

		if err := db.Create(&admin).Error; err != nil {
			return fmt.Errorf("failed to seed admin user: %w", err)
		}

		log.Println("✅ Seeded admin user: admin@inventory.com / admin123")
	}

	return nil
}
