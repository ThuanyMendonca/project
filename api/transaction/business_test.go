package transaction_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/ThuanyMendonca/project/api/transaction"
	"github.com/ThuanyMendonca/project/config/dbTransaction"
	"github.com/ThuanyMendonca/project/model"
	"github.com/ThuanyMendonca/project/repository"
	"github.com/ThuanyMendonca/project/service/authorization"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type TransactionBusinessSpy struct {
	DbTransactionFunc    dbTransaction.IDbTransaction
	TransactionRepo      repository.ITransactionRepository
	UserRepo             repository.IUserRepository
	UserType             repository.ITypeRepository
	BalanceRepo          repository.IBalanceRepository
	AuthorizationService authorization.IAuthorizationService
}

func TestGetTransactionByUserId(t *testing.T) {
	resp := &[]model.Transaction{}
	gofakeit.Slice(resp)

	tests := []struct {
		description         string
		userId              int64
		statusCodeExpected  int
		responseExpected    bool
		transactionResponse *[]model.Transaction
		erroExpected        error
		repositories        TransactionBusinessSpy
	}{
		{
			description:        "Should return error when userId is zero",
			statusCodeExpected: http.StatusBadRequest,
			userId:             0,
			erroExpected:       errors.New("não foi possível buscar transações para o usuário informado"),
		},
		{
			description:        "Should return not found to get transaction by user",
			statusCodeExpected: http.StatusNotFound,
			userId:             12,
			erroExpected:       errors.New("transações não encontradas para esse usuário"),
			repositories: TransactionBusinessSpy{
				TransactionRepo: &repository.TransactionRepositorySpy{
					GetByUserErr: gorm.ErrRecordNotFound,
				},
			},
		},
		{
			description:        "Should return not found to get transaction by user",
			statusCodeExpected: http.StatusInternalServerError,
			userId:             12,
			erroExpected:       errors.New("ocorreu um erro para buscar as transações do usuário"),
			repositories: TransactionBusinessSpy{
				TransactionRepo: &repository.TransactionRepositorySpy{
					GetByUserErr: errors.New("ocorreu um erro para buscar as transações do usuário"),
				},
			},
		},
		{
			description:        "Should return success to inactive user",
			statusCodeExpected: http.StatusOK,
			userId:             12,
			responseExpected:   true,
			repositories: TransactionBusinessSpy{
				TransactionRepo: &repository.TransactionRepositorySpy{
					GetByUserResp: resp,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("GetTransactionByUserId - %s", test.description), func(t *testing.T) {
			transactionNew := transaction.NewTransactionBusiness(test.repositories.DbTransactionFunc, test.repositories.TransactionRepo, test.repositories.UserRepo, test.repositories.UserType, test.repositories.BalanceRepo, test.repositories.AuthorizationService)

			statusCode, transactionResp, err := transactionNew.GetTransactionByUserId(test.userId)

			assert.Equal(t, test.statusCodeExpected, statusCode)

			if test.responseExpected {
				assert.NotEmpty(t, test.responseExpected, transactionResp)
			} else {
				assert.Empty(t, test.responseExpected, transactionResp)
			}

			assert.Equal(t, test.erroExpected, err)
		})
	}
}

func TestCreate(t *testing.T) {

	transactionResp := &model.Transaction{}
	gofakeit.Struct(transactionResp)

	userResp := &model.User{}
	gofakeit.Struct(userResp)
	userResp.IsActive = true

	balanceResp := &model.Balance{}
	gofakeit.Struct(balanceResp)
	balanceResp.Amount = 0

	tests := []struct {
		description        string
		transactionModel   *model.Transaction
		statusCodeExpected int
		erroExpected       error
		repositories       TransactionBusinessSpy
	}{
		{
			description:        "Should return error to get user by id",
			statusCodeExpected: http.StatusInternalServerError,
			transactionModel:   transactionResp,
			erroExpected:       errors.New("failed to get user"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdErr: errors.New("failed to get user"),
				},
			},
		},
		{
			description:        "Should return error to user inactive",
			statusCodeExpected: http.StatusBadRequest,
			transactionModel:   transactionResp,
			erroExpected:       errors.New("o usuário inativo, não é possível realizar transações"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: &model.User{
						IsActive: false,
					},
					GetByIdErr: nil,
				},
			},
		},
		{
			description:        "Should return error to get type",
			statusCodeExpected: http.StatusInternalServerError,
			transactionModel:   transactionResp,
			erroExpected:       errors.New("falha para pegar o tipo"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeErr:  errors.New("falha para pegar o tipo"),
				},
			},
		},
		{
			description:        "Should return error when type is shopkeeper",
			statusCodeExpected: http.StatusBadRequest,
			transactionModel:   transactionResp,
			erroExpected:       errors.New("lojista não pode realizar transações"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "SHOPKEEPER",
				},
			},
		},
		{
			description:        "Should return error not found balance",
			statusCodeExpected: http.StatusNotFound,
			transactionModel:   transactionResp,
			erroExpected:       errors.New("o usuário não possui saldo"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "COMMON",
				},
				BalanceRepo: &repository.BalanceRepositorySpy{
					GetBalance: nil,
					GetErr:     gorm.ErrRecordNotFound,
				},
			},
		},
		{
			description:        "Should return error to get balance",
			statusCodeExpected: http.StatusInternalServerError,
			transactionModel:   transactionResp,
			erroExpected:       errors.New("falha para pegar o saldo"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "COMMON",
				},
				BalanceRepo: &repository.BalanceRepositorySpy{
					GetBalance: nil,
					GetErr:     errors.New("falha para pegar o saldo"),
				},
			},
		},
		{
			description:        "Should return error when amount is zero",
			statusCodeExpected: http.StatusBadRequest,
			transactionModel:   transactionResp,
			erroExpected:       errors.New("saldo insuficiente"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "COMMON",
				},
				BalanceRepo: &repository.BalanceRepositorySpy{
					GetBalance: balanceResp,
				},
			},
		},
		{
			description:        "Should return error to begin transaction",
			statusCodeExpected: http.StatusInternalServerError,
			transactionModel:   transactionResp,
			erroExpected:       errors.New("falha para abrir uma transaction no banco de dados"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "COMMON",
				},
				BalanceRepo: &repository.BalanceRepositorySpy{
					GetBalance: &model.Balance{
						Amount: 12,
					},
				},
				DbTransactionFunc: &dbTransaction.DbTransactionSpy{
					BeginErr: errors.New("falha para abrir uma transaction no banco de dados"),
				},
			},
		},
		{
			description:        "Should return error to discount amount",
			statusCodeExpected: http.StatusBadRequest,
			transactionModel:   transactionResp,
			erroExpected:       errors.New("saldo insuficiente para realizar a transação"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "COMMON",
				},
				BalanceRepo: &repository.BalanceRepositorySpy{
					GetBalance: &model.Balance{
						Amount: 1000,
					},
				},
				DbTransactionFunc: &dbTransaction.DbTransactionSpy{},
			},
		},
		{
			description:        "Should return error to request autorization",
			statusCodeExpected: http.StatusInternalServerError,
			transactionModel: &model.Transaction{
				Value: 20,
			},
			erroExpected: errors.New("ocorreu um erro para solicitar autorização da transação"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "COMMON",
				},
				BalanceRepo: &repository.BalanceRepositorySpy{
					GetBalance: &model.Balance{
						Amount: 10000,
					},
				},
				DbTransactionFunc: &dbTransaction.DbTransactionSpy{},
				AuthorizationService: &authorization.AuthorizationServiceSpy{
					StatusCode:        http.StatusInternalServerError,
					AuthorizationErr:  errors.New("ocorreu um erro para solicitar autorização da transação"),
					AuthorizationResp: nil,
				},
			},
		},
		{
			description:        "Should return transaction not autorized",
			statusCodeExpected: http.StatusOK,
			transactionModel: &model.Transaction{
				Value: 20,
			},
			erroExpected: errors.New("transação não autorizada"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "COMMON",
				},
				BalanceRepo: &repository.BalanceRepositorySpy{
					GetBalance: &model.Balance{
						Amount: 10000,
					},
				},
				DbTransactionFunc: &dbTransaction.DbTransactionSpy{},
				AuthorizationService: &authorization.AuthorizationServiceSpy{
					StatusCode: http.StatusOK,
					AuthorizationResp: &model.AuthorizeResponse{
						Authorization: false,
					},
				},
			},
		},
		{
			description:        "Should return error to insert transaction in db",
			statusCodeExpected: http.StatusInternalServerError,
			transactionModel: &model.Transaction{
				Value: 20,
			},
			erroExpected: errors.New("falha ao inserir a transação no banco"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "COMMON",
				},
				BalanceRepo: &repository.BalanceRepositorySpy{
					GetBalance: &model.Balance{
						Amount: 10000,
					},
				},
				DbTransactionFunc: &dbTransaction.DbTransactionSpy{},
				AuthorizationService: &authorization.AuthorizationServiceSpy{
					StatusCode: http.StatusOK,
					AuthorizationResp: &model.AuthorizeResponse{
						Authorization: true,
					},
				},
				TransactionRepo: &repository.TransactionRepositorySpy{
					TransactionErr: errors.New("falha ao inserir a transação no banco"),
				},
			},
		},
		{
			description:        "Should return error to update balance",
			statusCodeExpected: http.StatusInternalServerError,
			transactionModel: &model.Transaction{
				Value: 20,
			},
			erroExpected: errors.New("ocorreu erro ao atualizar o saldo"),
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "COMMON",
				},
				BalanceRepo: &repository.BalanceRepositorySpy{
					GetBalance: &model.Balance{
						Amount: 10000,
					},
					UpdateValueErr: errors.New("ocorreu erro ao atualizar o saldo"),
				},
				DbTransactionFunc: &dbTransaction.DbTransactionSpy{},
				AuthorizationService: &authorization.AuthorizationServiceSpy{
					StatusCode: http.StatusOK,
					AuthorizationResp: &model.AuthorizeResponse{
						Authorization: true,
					},
				},
				TransactionRepo: &repository.TransactionRepositorySpy{},
			},
		},
		{
			description:        "Should return success",
			statusCodeExpected: http.StatusCreated,
			transactionModel: &model.Transaction{
				Value: 20,
			},
			repositories: TransactionBusinessSpy{
				UserRepo: &repository.UserRepositorySpy{
					GetByIdResp: userResp,
					GetTypeResp: "COMMON",
				},
				BalanceRepo: &repository.BalanceRepositorySpy{
					GetBalance: &model.Balance{
						Amount: 10000,
					},
				},
				DbTransactionFunc: &dbTransaction.DbTransactionSpy{},
				AuthorizationService: &authorization.AuthorizationServiceSpy{
					StatusCode: http.StatusOK,
					AuthorizationResp: &model.AuthorizeResponse{
						Authorization: true,
					},
				},
				TransactionRepo: &repository.TransactionRepositorySpy{},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Create - %s", test.description), func(t *testing.T) {
			transactionNew := transaction.NewTransactionBusiness(test.repositories.DbTransactionFunc, test.repositories.TransactionRepo, test.repositories.UserRepo, test.repositories.UserType, test.repositories.BalanceRepo, test.repositories.AuthorizationService)

			statusCode, err := transactionNew.Create(test.transactionModel)

			assert.Equal(t, test.statusCodeExpected, statusCode)
			assert.Equal(t, test.erroExpected, err)
		})
	}
}
