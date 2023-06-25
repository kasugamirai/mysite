package prize_models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Code struct {
	gorm.Model
	Code   string `gorm:"unique"`
	IsUsed bool   // 是否已被使用
}

// 添加新的兑换码到数据库中
func AddCode(db *gorm.DB, code string) error {
	c := Code{
		Code:   code,
		IsUsed: false,
	}
	if err := db.Create(&c).Error; err != nil {
		return fmt.Errorf("failed to add code: %w", err)
	}
	return nil
}

// 从数据库中获取未使用的兑换码
func GetCode(db *gorm.DB) (string, error) {
	var code Code
	// 获取第一个未使用的兑换码
	if err := db.Where("is_used = ?", false).First(&code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("no codes left")
		}
		return "", err
	}

	// 将获取到的兑换码标记为已使用
	code.IsUsed = true
	if err := db.Save(&code).Error; err != nil {
		return "", fmt.Errorf("failed to mark code as used: %w", err)
	}

	return code.Code, nil
}
