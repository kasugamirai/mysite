package prize_models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type RedemptionCode struct {
	gorm.Model
	Code string `gorm:"unique"`
	Used bool   `gorm:"default:false"` // Add 'Used' field
}

func AddRedemptionCode(db *gorm.DB, code string) error {
	c := RedemptionCode{
		Code: code,
		Used: false,
	}
	if err := db.Create(&c).Error; err != nil {
		return fmt.Errorf("failed to add code: %w", err)
	}
	return nil
}

func UseRedemptionCode(db *gorm.DB, code string) error {
	var c RedemptionCode

	if err := db.Where("code = ?", code).First(&c).Error; err != nil {
		// 如果在数据库中找不到这个兑换码，返回错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("code not found")
		}
		// 如果在查找过程中出现其他错误，返回错误
		return fmt.Errorf("failed to find code: %w", err)
	}

	// 如果这个兑换码已经被使用，返回一个错误
	if c.Used {
		return fmt.Errorf("code already used")
	}

	// 如果找到了这个兑换码，并且还没有被使用，将其标记为已使用
	c.Used = true
	if err := db.Save(&c).Error; err != nil {
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	return nil
}
