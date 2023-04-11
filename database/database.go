package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"xy.com/mysite/config"
	"xy.com/mysite/models"
)

var (
	DB *gorm.DB
)

// InitDB initializes the database connection.
func InitDB() error {
	var err error
	dbConfig := config.Config{}

	// Create a new SQLite database connection.
	dsn := dbConfig.DatabaseDSN
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the data models.
	err = migrateModels()
	if err != nil {
		return err
	}

	return nil
}

// migrateModels migrates the data models to the database.
func migrateModels() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		return err
	}

	return nil
}
