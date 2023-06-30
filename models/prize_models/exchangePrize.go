package prize_models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ExchangedPrize struct {
	gorm.Model
	PrizeName      string `json:"prize_name"`
	UserID         string `json:"user_id"`
	RedemptionCode string `json:"redemption_code"`
}

// ExchangePrize exchanges a prize
func ExchangePrize(db *gorm.DB, userID string, prizeName string, pointsSystem *PointsSystem) (string, error) {
	// Check if the prize exists
	prize, err := GetPrizeByName(db, prizeName)
	if err != nil {
		return "", err
	}

	// Check if the user has already exchanged this prize
	hasExchanged, err := CheckIfUserExchangedPrize(db, userID, prizeName)
	if err != nil {
		return "", err
	}

	// If the user has already exchanged this prize, return the redemption code
	if hasExchanged {
		redemptionCode, err := GetRedemptionCode(db, userID, prizeName)
		if err != nil {
			return "", fmt.Errorf("failed to get redemption code: %w", err)
		}
		return redemptionCode, nil
	}

	// The user has not exchanged this prize yet and can exchange it
	if pointsSystem.Points < prize.Cost {
		// The user does not have enough points to exchange the prize
		return "", errors.New("insufficient points")
	}

	// Generate a redemption code
	redemptionCode, err := GetCode(db)
	if err != nil {
		return "", fmt.Errorf("failed to get redemption code: %w", err)
	}

	// Create an ExchangedPrize record
	exchangedPrize := ExchangedPrize{
		PrizeName:      prizeName,
		UserID:         userID,
		RedemptionCode: redemptionCode,
	}

	// Save the ExchangedPrize to the database
	if err := db.Create(&exchangedPrize).Error; err != nil {
		return "", fmt.Errorf("failed to save exchanged prize: %w", err)
	}

	return redemptionCode, nil
}

// CheckIfUserExchangedPrize checks if a user has already exchanged a prize
func CheckIfUserExchangedPrize(db *gorm.DB, userID string, prizeName string) (bool, error) {
	var exchangedPrize ExchangedPrize
	err := db.Where("user_id = ? AND prize_name = ?", userID, prizeName).First(&exchangedPrize).Error

	// If we get an error, and it's a ErrRecordNotFound, we return false (has not exchanged), and no error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		// If the error is something else, we return it
		return false, fmt.Errorf("failed to query prize_handlers: %w", err)
	}

	// If we found a record, the user has already exchanged this prize
	return true, nil
}

func GetPrizeByName(db *gorm.DB, prizeName string) (*Prize, error) {
	var prize Prize
	err := db.Where("prize_name = ?", prizeName).First(&prize).Error

	if err != nil {
		// We didn't find a record with the given prize name
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("prize_handlers not found: %s", prizeName)
		}

		// Some other error occurred
		return nil, fmt.Errorf("failed to query prize_handlers: %w", err)
	}

	// We found a record, so we return it
	return &prize, nil
}

// GetRedemptionCode fetches the redemption code for a given user and prize
func GetRedemptionCode(db *gorm.DB, userID string, prizeName string) (string, error) {
	var exchangedPrize ExchangedPrize
	err := db.Where("user_id = ? AND prize_name = ?", userID, prizeName).First(&exchangedPrize).Error

	if err != nil {
		// We didn't find a record with the given user ID and prize name
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("exchanged prize not found for user %s and prize %s", userID, prizeName)
		}

		// Some other error occurred
		return "", fmt.Errorf("failed to query exchanged prizes: %w", err)
	}

	// We found a record, so we return the redemption code
	return exchangedPrize.RedemptionCode, nil
}
