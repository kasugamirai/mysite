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

// ExchangePrize 兑换奖品
func ExchangePrize(db *gorm.DB, userID string, prizeName string, pointsSystem *PointsSystem) (string, error) {
	prize, err := GetPrizeByName(db, prizeName)
	if err != nil {
		return "", err
	}

	hasExchanged, err := CheckIfUserExchangedPrize(db, userID, prizeName)
	if err != nil {
		return "", err
	}
	if hasExchanged {
		// 用户已经兑换过这个奖品
		return "", fmt.Errorf("user %s has already exchanged prize %s", userID, prizeName)
	}

	// 用户还没有兑换过这个奖品，可以兑换
	if pointsSystem.Points < prize.Cost {
		return "", errors.New("insufficient points")
	}
	// 不再减少积分
	// pointsSystem.Points -= prize.Cost

	redemptionCode, err := GetCode(db)
	if err != nil {
		return "", fmt.Errorf("failed to get redemption code: %w", err)
	}

	exchangedPrize := ExchangedPrize{
		PrizeName:      prizeName,
		UserID:         userID,
		RedemptionCode: redemptionCode,
	}
	if err := db.Create(&exchangedPrize).Error; err != nil {
		return "", fmt.Errorf("failed to save exchanged prize: %w", err)
	}
	return redemptionCode, nil
}

// CheckIfUserExchangedPrize 检查用户是否已经兑换过这个奖品
func CheckIfUserExchangedPrize(db *gorm.DB, userID string, prizeName string) (bool, error) {
	var exchangedPrize ExchangedPrize
	if err := db.Where("user_id = ? AND prize_name = ?", userID, prizeName).First(&exchangedPrize).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户还没有兑换过这个奖品
			return false, nil
		}
		// 数据库错误
		return false, fmt.Errorf("failed to query prize_handlers: %w", err)
	}
	// 用户已经兑换过这个奖品
	return true, nil
}

func GetPrizeByName(db *gorm.DB, prizeName string) (*Prize, error) {
	var prize Prize
	if err := db.Where("prize_name = ?", prizeName).First(&prize).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("prize_handlers not found: %s", prizeName)
		}
		return nil, fmt.Errorf("failed to query prize_handlers: %w", err)
	}
	return &prize, nil
}
