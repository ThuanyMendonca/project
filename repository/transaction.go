package repository

import (
	"github.com/ThuanyMendonca/project/model"
	"gorm.io/gorm"
)

type ITransactionRepository interface {
	Create(transaction *model.Transaction) error
	GetByUser(userId int64) (*[]model.Transaction, error)
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionRepository{db}
}

func (t *TransactionRepository) Create(transaction *model.Transaction) error {
	return t.db.Create(transaction).Error
}

func (t *TransactionRepository) GetByUser(userId int64) (*[]model.Transaction, error) {
	transaction := &[]model.Transaction{}
	find := t.db

	if userId != 0 {
		find = find.Where("userId = ?", userId)
	}

	find = find.Joins("User")
	err := find.Find(transaction).Error

	return transaction, err

}
