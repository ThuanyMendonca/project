package transaction_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThuanyMendonca/project/api/transaction"
	"github.com/ThuanyMendonca/project/model"
	"github.com/ThuanyMendonca/project/utils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TransactionControllerSpy struct {
	business transaction.ITransactionBusiness
}

func TestGetTransactionByUserIdController(t *testing.T) {
	resp := &[]model.Transaction{}
	gofakeit.Slice(resp)

	transaction.Router(&gin.Default().RouterGroup)
	tests := []struct {
		description        string
		userId             string
		key                string
		payload            *bytes.Buffer
		codeStatusExpected int
		errorExpected      string
		services           TransactionControllerSpy
	}{
		{
			description:        "Should return an error if not sent userId",
			codeStatusExpected: http.StatusBadRequest,
			errorExpected:      "id do usuário é obrigatório",
			payload:            utils.MockValidPayload(resp),
		},
		{
			description:        "Should return an error to get transactions",
			userId:             "123",
			key:                "userId",
			codeStatusExpected: http.StatusInternalServerError,
			errorExpected:      "ocorreu um erro para pegar as transações do usuário",
			payload:            utils.MockValidPayload(resp),
			services: TransactionControllerSpy{
				business: &transaction.TransactionBusinessSpy{
					GetTransactionByUserIdErr:        errors.New("ocorreu um erro para pegar as transações do usuário"),
					GetTransactionByUserIdStatusCode: http.StatusInternalServerError,
				},
			},
		},
		{
			description:        "Should return success to get transactions",
			userId:             "123",
			key:                "userId",
			codeStatusExpected: http.StatusOK,
			payload:            utils.MockValidPayload(resp),
			services: TransactionControllerSpy{
				business: &transaction.TransactionBusinessSpy{
					GetTransactionByUserIdResp:       resp,
					GetTransactionByUserIdStatusCode: http.StatusOK,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("GetTransactionByUserId - %s", test.description), func(t *testing.T) {
			c := transaction.NewUTransactionController(test.services.business)

			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = append(ctx.Params, gin.Param{
				Key:   test.key,
				Value: test.userId,
			})

			ctx.Request = httptest.NewRequest(http.MethodGet, "/transaction/", test.payload)

			c.GetTransactionByUserId(ctx)

			if test.errorExpected != "" {
				resp := map[string]string{}
				byte, _ := ioutil.ReadAll(w.Result().Body)
				if err := json.Unmarshal(byte, &resp); err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.codeStatusExpected, w.Result().StatusCode)
				assert.Equal(t, test.errorExpected, resp["message"])
			}
		})
	}
}

func TestPostController(t *testing.T) {
	response := &model.Transaction{}
	gofakeit.Slice(response)

	transaction.Router(&gin.Default().RouterGroup)

	tests := []struct {
		description        string
		payload            *bytes.Buffer
		codeStatusExpected int
		errorExpected      string
		services           TransactionControllerSpy
	}{
		{
			description:        "Should return an error if sent invalid payload",
			codeStatusExpected: http.StatusBadRequest,
			errorExpected:      "conteúdo da requisição inválido",
			payload:            utils.MockInvalidPayload(),
		},
		{
			description:        "Should return an error to validate body transactions",
			codeStatusExpected: http.StatusBadRequest,
			errorExpected:      "valor é obrigatório",
			payload: utils.MockValidPayload(&model.Transaction{
				Value: 0,
			}),
		},
		{
			description:        "Should return an error to create transactions",
			codeStatusExpected: http.StatusInternalServerError,
			errorExpected:      "ocorreu um erro para criar a transação do usuário",
			payload:            utils.MockValidPayload(response),
			services: TransactionControllerSpy{
				business: &transaction.TransactionBusinessSpy{
					CreateStatusCode:     http.StatusInternalServerError,
					CreateTransactionErr: errors.New("ocorreu um erro para criar a transação do usuário"),
				},
			},
		},
		{
			description:        "Should return success to get transactions",
			codeStatusExpected: http.StatusOK,
			payload:            utils.MockValidPayload(response),
			services: TransactionControllerSpy{
				business: &transaction.TransactionBusinessSpy{
					CreateStatusCode:     http.StatusCreated,
					CreateTransactionErr: nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Post - %s", test.description), func(t *testing.T) {
			c := transaction.NewUTransactionController(test.services.business)

			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = httptest.NewRequest(http.MethodPost, "/transaction", test.payload)

			c.Post(ctx)

			if test.errorExpected != "" {
				resp := map[string]string{}
				byte, _ := ioutil.ReadAll(w.Result().Body)
				if err := json.Unmarshal(byte, &resp); err != nil {
					t.Error(err)
				}

				assert.Equal(t, test.codeStatusExpected, w.Result().StatusCode)
				assert.Equal(t, test.errorExpected, resp["message"])
			}
		})
	}
}
