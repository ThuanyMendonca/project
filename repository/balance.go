package repository

import (
	"time"

	"github.com/ThuanyMendonca/project/model"
	"gorm.io/gorm"
)

type IBalanceRepository interface {
	Add(balance *model.Balance) error
	Get(userId int64) (*model.Balance, error)
	UpdateValue(db *gorm.DB, userId int64, newValue float64) error
	Update(userId int64, newValue float64) error
}

type BalanceReposity struct {
	db *gorm.DB
}

func NewBalanceRepository(db *gorm.DB) IBalanceRepository {
	return &BalanceReposity{db}
}

func (b *BalanceReposity) Add(balance *model.Balance) error {
	return b.db.Create(balance).Error
}

func (b *BalanceReposity) Get(userId int64) (*model.Balance, error) {
	balance := &model.Balance{}
	find := b.db

	if userId != 0 {
		find = find.Where("balances.user_id = ?", userId)
	}

	find = find.Joins("User")
	err := find.First(balance).Error

	if err != nil {
		return nil, err
	}

	return balance, nil
}

// update para atualizar o amount caso tenha quando fizer a transação
func (b *BalanceReposity) UpdateValue(db *gorm.DB, userId int64, newValue float64) error {
	return db.Model(&model.Balance{}).Where("user_id = ?", userId).UpdateColumns(map[string]interface{}{
		"value":      newValue,
		"updated_at": time.Now(),
	}).Error
}

func (b *BalanceReposity) Update(userId int64, newValue float64) error {
	return b.db.Model(&model.Balance{}).Where("user_id = ?", userId).UpdateColumns(map[string]interface{}{
		"value":      newValue,
		"updated_at": time.Now(),
	}).Error
}
