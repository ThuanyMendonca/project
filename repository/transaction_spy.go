package repository

import (
	"github.com/ThuanyMendonca/project/model"
	"gorm.io/gorm"
)

type TransactionRepositorySpy struct {
	ITransactionRepository
	TransactionErr error
	GetByUserResp  *[]model.Transaction
	GetByUserErr   error
}

func (t *TransactionRepositorySpy) Create(tx *gorm.DB, transaction *model.Transaction) error {
	return t.TransactionErr
}

func (t *TransactionRepositorySpy) GetByUser(userId int64) (*[]model.Transaction, error) {
	return t.GetByUserResp, t.GetByUserErr
}
