package repository

import (
	"github.com/ThuanyMendonca/project/model"
	"gorm.io/gorm"
)

type BalanceRepositorySpy struct {
	IBalanceRepository
	AddErr         error
	GetBalance     *model.Balance
	GetErr         error
	UpdateValueErr error
	UpdateErr      error
}

func (b *BalanceRepositorySpy) Add(balance *model.Balance) error {
	return b.AddErr
}

func (b *BalanceRepositorySpy) Get(userId int64) (*model.Balance, error) {
	return b.GetBalance, b.GetErr
}

func (b *BalanceRepositorySpy) UpdateValue(db *gorm.DB, userId int64, newValue float64) error {
	return b.UpdateValueErr
}

func (b *BalanceRepositorySpy) Update(userId int64, newValue float64) error {
	return b.UpdateErr
}
