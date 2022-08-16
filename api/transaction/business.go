package transaction

import (
	"errors"
	"net/http"

	"github.com/ThuanyMendonca/project/model"
	"github.com/ThuanyMendonca/project/repository"
	"gorm.io/gorm"
)

type ITransaction interface {
}

type TransactionBusiness struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
}

func NewTransactionBusiness(transactionRepo repository.TransactionRepository, userRepo repository.UserRepository, userType repository.TypeRepository) ITransaction {
	return &TransactionBusiness{transactionRepo, userRepo}
}

func (t *TransactionBusiness) GetTransactionByUserId(userId int64) (int, *[]model.Transaction, error) {
	if userId == 0 {
		return http.StatusBadRequest, nil, errors.New("não foi possível buscar transações para o usuário informado")
	}

	transactions, err := t.transactionRepo.GetByUser(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, nil, errors.New("transações não encontradas para esse usuário")
		}
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, transactions, nil
}

func (t *TransactionBusiness) Create(transaction *model.Transaction) (int, error) {
	isActiveUser, err := t.userRepo.IsActive(transaction.PayerId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !isActiveUser {
		return http.StatusBadRequest, errors.New("o usuário inativo, não é possível realizar transações")
	}

	userType, err := t.userRepo.GetType(transaction.PayerId)
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, err
	}

	shopper := "SHOPKEEPER"
	if userType == &shopper {
		return http.StatusBadRequest, errors.New("lojista não pode realizar transações")
	}

	// verificar se tem saldo

	return 0, nil
}
