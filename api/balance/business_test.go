package balance_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/ThuanyMendonca/project/api/balance"
	"github.com/ThuanyMendonca/project/model"
	"github.com/ThuanyMendonca/project/repository"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type BalanceBusinessRepositorySpy struct {
	BalanceRepository repository.IBalanceRepository
}

func TestCreate(t *testing.T) {

	getBalanceResp := &model.Balance{}
	gofakeit.Struct(getBalanceResp)

	tests := []struct {
		description        string
		userId             int64
		value              float64
		statusCodeExpected int
		erroExpected       error
		repositories       BalanceBusinessRepositorySpy
	}{
		{
			description:        "Should return error to get user to create balance",
			statusCodeExpected: http.StatusInternalServerError,
			erroExpected:       errors.New("falha para pegar o saldo"),
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: nil,
					GetErr:     errors.New("falha para pegar o saldo"),
				},
			},
		},
		{
			description:        "Should return balance is exists",
			statusCodeExpected: http.StatusBadRequest,
			erroExpected:       errors.New("já existe um registro de saldo, considere atualizar o saldo"),
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: getBalanceResp,
					GetErr:     nil,
				},
			},
		},
		{
			description:        "Should return error to create balance in db",
			statusCodeExpected: http.StatusInternalServerError,
			erroExpected:       errors.New("falha para inserir o saldo"),
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					AddErr: errors.New("falha para inserir o saldo"),
				},
			},
		},
		{
			description:        "Should return success to create balance",
			statusCodeExpected: http.StatusOK,
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: nil,
					GetErr:     nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Create - %s", test.description), func(t *testing.T) {
			b := balance.NewBalanceBusiness(test.repositories.BalanceRepository)

			statusCode, err := b.Create(test.userId, test.value)

			assert.Equal(t, test.statusCodeExpected, statusCode)
			assert.Equal(t, test.erroExpected, err)
		})
	}
}

func TestUpdate(t *testing.T) {

	getBalanceResp := &model.Balance{}
	gofakeit.Struct(getBalanceResp)

	tests := []struct {
		description        string
		userId             int64
		value              float64
		statusCodeExpected int
		erroExpected       error
		repositories       BalanceBusinessRepositorySpy
	}{
		{
			description:        "Should return error to get user to update balance",
			statusCodeExpected: http.StatusInternalServerError,
			erroExpected:       errors.New("failed to get balance"),
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: nil,
					GetErr:     errors.New("failed to get balance"),
				},
			},
		},
		{
			description:        "Should get user not found",
			statusCodeExpected: http.StatusNotFound,
			erroExpected:       errors.New("nenhum registro encontrado para ser atualizado"),
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: nil,
					GetErr:     gorm.ErrRecordNotFound,
				},
			},
		},
		{
			description:        "Should get user not found",
			statusCodeExpected: http.StatusInternalServerError,
			erroExpected:       errors.New("failed get"),
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: nil,
					GetErr:     errors.New("failed get"),
				},
			},
		},
		{
			description:        "Should return error to calculate new balance",
			statusCodeExpected: http.StatusInternalServerError,
			erroExpected:       errors.New("não é possível adicionar saldo com valor zerado"),
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: getBalanceResp,
					GetErr:     nil,
					UpdateErr:  errors.New("não é possível adicionar saldo com valor zerado"),
				},
			},
		},
		{
			description:        "Should return error to update balance",
			statusCodeExpected: http.StatusInternalServerError,
			erroExpected:       errors.New("ocorreu um erro ao atualizar o saldo"),
			userId:             11,
			value:              12,
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: getBalanceResp,
					GetErr:     nil,
					UpdateErr:  errors.New("ocorreu um erro ao atualizar o saldo"),
				},
			},
		},
		{
			description:        "Should return success to create balance",
			statusCodeExpected: http.StatusOK,
			userId:             123,
			value:              12.4,
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: getBalanceResp,
					GetErr:     nil,
					UpdateErr:  nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Update - %s", test.description), func(t *testing.T) {
			b := balance.NewBalanceBusiness(test.repositories.BalanceRepository)

			statusCode, err := b.Update(test.userId, test.value)

			assert.Equal(t, test.statusCodeExpected, statusCode)
			assert.Equal(t, test.erroExpected, err)
		})
	}
}

func TestGet(t *testing.T) {
	getBalanceResp := &model.Balance{}
	gofakeit.Struct(getBalanceResp)

	tests := []struct {
		description        string
		userId             int64
		statusCodeExpected int
		responseExpected   bool
		erroExpected       error
		repositories       BalanceBusinessRepositorySpy
	}{
		{
			description:        "Should return error to get user",
			statusCodeExpected: http.StatusInternalServerError,
			erroExpected:       errors.New("falha para buscar o saldo"),
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: nil,
					GetErr:     errors.New("falha para buscar o saldo"),
				},
			},
		},
		{
			description:        "Should get user not found",
			statusCodeExpected: http.StatusNotFound,
			erroExpected:       errors.New("registro não encontrado"),
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: nil,
					GetErr:     gorm.ErrRecordNotFound,
				},
			},
		},
		{
			description:        "Should return success to get balance",
			statusCodeExpected: http.StatusOK,
			responseExpected:   true,
			userId:             123,
			repositories: BalanceBusinessRepositorySpy{
				BalanceRepository: &repository.BalanceRepositorySpy{
					GetBalance: getBalanceResp,
					GetErr:     nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Get - %s", test.description), func(t *testing.T) {
			b := balance.NewBalanceBusiness(test.repositories.BalanceRepository)

			statusCode, balance, err := b.Get(test.userId)

			assert.Equal(t, test.statusCodeExpected, statusCode)

			if test.responseExpected {
				assert.NotEmpty(t, test.responseExpected, balance)
			} else {
				assert.Empty(t, test.responseExpected, balance)
			}

			assert.Equal(t, test.erroExpected, err)
		})
	}
}
