package transaction

import "github.com/ThuanyMendonca/project/model"

type TransactionBusinessSpy struct {
	ITransactionBusiness
	GetTransactionByUserIdStatusCode int
	GetTransactionByUserIdResp       *[]model.Transaction
	GetTransactionByUserIdErr        error
	CreateStatusCode                 int
	CreateTransactionErr             error
}

func (t *TransactionBusinessSpy) GetTransactionByUserId(userId int64) (int, *[]model.Transaction, error) {
	return t.GetTransactionByUserIdStatusCode, t.GetTransactionByUserIdResp, t.GetTransactionByUserIdErr
}

func (t *TransactionBusinessSpy) Create(transaction *model.Transaction) (int, error) {
	return t.CreateStatusCode, t.CreateTransactionErr
}
