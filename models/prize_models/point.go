package prize_models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type PointsSystem struct {
	gorm.Model
	UserID string `json:"user_id"`
	Points int    `json:"points"`
	Coins  int    `json:"coins"`
}

// 通过用户ID获取积分系统
func GetPointsSystem(db *gorm.DB, userID string) (*PointsSystem, error) {
	var pointsSystem PointsSystem
	if err := db.Where("user_id = ?", userID).First(&pointsSystem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User does not exist, create a new points system with zero points
			pointsSystem = PointsSystem{UserID: userID, Points: 0, Coins: 0}
			if err := db.Create(&pointsSystem).Error; err != nil {
				return nil, fmt.Errorf("failed to create points system: %w", err)
			}
		} else {
			return nil, err
		}
	}
	return &pointsSystem, nil
}

// 更新积分系统
func UpdatePointsSystem(db *gorm.DB, pointsSystem *PointsSystem) error {
	if err := db.Save(pointsSystem).Error; err != nil {
		return fmt.Errorf("failed to update points system: %w", err)
	}
	return nil
}

func (ps *PointsSystem) Draw(db *gorm.DB) error {
	if ps.Points < 40000 {
		// 积分低于40000时，100%概率获得1000积分
		ps.Points += 1000
	} else if ps.Points >= 40000 && ps.Points < 50000 {
		// 积分高于40000时，每次抽奖获得金币
		ps.Coins += 10
	} else {
		return errors.New("积分已经超过50000，不能再抽奖")
	}

	// 更新积分系统到数据库
	if err := UpdatePointsSystem(db, ps); err != nil {
		return err
	}

	return nil
}

func (ps *PointsSystem) ExchangeCoins(db *gorm.DB) error {
	if ps.Coins < 100 {
		// 金币不足100，不能兑换
		return errors.New("金币不足，不能兑换")
	}

	// 每100金币可以兑换1积分
	ps.Points += ps.Coins / 100
	ps.Coins %= 100

	// 更新积分系统到数据库
	if err := UpdatePointsSystem(db, ps); err != nil {
		return err
	}

	return nil
}
