package transaction

import (
	"errors"
	"net/http"
	"time"

	"github.com/ThuanyMendonca/project/config/dbTransaction"
	"github.com/ThuanyMendonca/project/model"
	"github.com/ThuanyMendonca/project/repository"
	"github.com/ThuanyMendonca/project/service/authorization"
	"gorm.io/gorm"
)

type ITransactionBusiness interface {
	Create(transaction *model.Transaction) (int, error)
	GetTransactionByUserId(userId int64) (int, *[]model.Transaction, error)
}

type TransactionBusiness struct {
	dbTransactionFunc    dbTransaction.IDbTransaction
	transactionRepo      repository.ITransactionRepository
	userRepo             repository.IUserRepository
	userType             repository.ITypeRepository
	balanceRepo          repository.IBalanceRepository
	authorizationService authorization.IAuthorizationService
}

func NewTransactionBusiness(dbTransactionFunc dbTransaction.IDbTransaction, transactionRepo repository.ITransactionRepository, userRepo repository.IUserRepository, userType repository.ITypeRepository, balanceRepo repository.IBalanceRepository, authorizationService authorization.IAuthorizationService) ITransactionBusiness {
	return &TransactionBusiness{dbTransactionFunc, transactionRepo, userRepo, userType, balanceRepo, authorizationService}
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
	user, err := t.userRepo.GetById(transaction.PayerId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !user.IsActive {
		return http.StatusBadRequest, errors.New("o usuário inativo, não é possível realizar transações")
	}

	// Verifica o tipo de usuário
	userType, err := t.userRepo.GetType(user.TypeId)
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, err
	}

	shopper := "SHOPKEEPER"
	if userType == shopper {
		return http.StatusBadRequest, errors.New("lojista não pode realizar transações")
	}

	// verificar se tem saldo
	balance, err := t.balanceRepo.Get(transaction.PayerId)
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, errors.New("o usuário não possui saldo")
	}

	if balance.Amount <= 0 {
		return http.StatusBadRequest, errors.New("saldo insuficiente")
	}

	tx, err := t.dbTransactionFunc.Begin()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	newAmount, err := t.discountAmount(balance.Amount, transaction.Value)
	if err != nil {
		t.dbTransactionFunc.Rollback(tx)
		return http.StatusBadRequest, err
	}

	statusCode, authorizated, err := t.authorizationService.Authorize()
	if err != nil {
		t.dbTransactionFunc.Rollback(tx)
		return statusCode, err
	}

	if !authorizated.Authorization {
		t.dbTransactionFunc.Rollback(tx)
		return http.StatusInternalServerError, errors.New("transação não autorizada")
	}

	// Cria a transação
	addTransaction := model.Transaction{
		Value:     transaction.Value,
		PayerId:   transaction.PayerId,
		CreatedAt: time.Now(),
	}

	if err := t.createTransaction(tx, addTransaction); err != nil {
		t.dbTransactionFunc.Rollback(tx)
		return http.StatusInternalServerError, err
	}

	// atualizar o saldo na tabela balance
	if err := t.balanceRepo.UpdateValue(tx, transaction.PayerId, *newAmount); err != nil {
		t.dbTransactionFunc.Rollback(tx)
		return http.StatusInternalServerError, err
	}

	t.dbTransactionFunc.Commit(tx)
	return http.StatusCreated, nil
}

func (t *TransactionBusiness) discountAmount(balanceAmount, transactionAmount float64) (*float64, error) {
	calculate := balanceAmount - transactionAmount

	if calculate < 0 {
		return nil, errors.New("saldo insuficiente para realizar a transação")
	}

	return &calculate, nil
}

func (t *TransactionBusiness) createTransaction(tx *gorm.DB, transaction model.Transaction) error {

	if err := t.transactionRepo.Create(tx, &transaction); err != nil {
		return err
	}
	return nil
}
