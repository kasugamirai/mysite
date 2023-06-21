package shop_models

import (
	"gorm.io/gorm"
)

// Product represents a product entity in the system.
type Product struct {
	gorm.Model
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"size:1024" json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	ImageURL    string  `gorm:"size:512" json:"image_url"`
}

// CreateProduct creates a new product in the database.
func CreateProduct(db *gorm.DB, product *Product) error {
	return db.Create(product).Error
}

// GetProductByID retrieves a product from the database by ID.
func GetProductByID(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	err := db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAllProducts retrieves all products from the database.
func GetAllProducts(db *gorm.DB) ([]Product, error) {
	var products []Product
	err := db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// UpdateProduct updates the product data in the database.
func UpdateProduct(db *gorm.DB, product *Product) error {
	return db.Save(product).Error
}

// DeleteProduct deletes a product from the database.
func DeleteProduct(db *gorm.DB, id uint) error {
	return db.Delete(&Product{}, id).Error
}
