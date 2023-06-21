package prize_models

import (
	"fmt"
	"gorm.io/gorm"
)

type Prize struct {
	gorm.Model
	PrizeName string `json:"prize_name"`
	Cost      int    `json:"cost"` // 积分兑换的价格
}

// 添加新的奖品到数据库中
func AddPrize(db *gorm.DB, prizeName string, cost int) error {
	prize := Prize{
		PrizeName: prizeName,
		Cost:      cost,
	}
	if err := db.Create(&prize).Error; err != nil {
		return fmt.Errorf("failed to add prize_handlers: %w", err)
	}
	return nil
}
