package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"xy.com/mysite/config"
	"xy.com/mysite/models/prize_models"
	"xy.com/mysite/models/shop_models"
	"xy.com/mysite/models/user_models"
)

var (
	DB *gorm.DB
)

// InitDB initializes the database connection.
func InitDB() error {
	var err error
	dbConfig := config.Instance.DatabaseDSN

	// Create a new SQLite database connection.
	DB, err = gorm.Open(sqlite.Open(dbConfig), &gorm.Config{})
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
		&user_models.User{},
		&shop_models.Product{},
		&shop_models.Order{},
		&shop_models.OrderItem{},
		&prize_models.Prize{},
		&prize_models.ExchangedPrize{},
		&prize_models.PointsSystem{},
	)
	if err != nil {
		return err
	}

	return nil
}
