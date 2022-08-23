package balance_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThuanyMendonca/project/api/balance"
	"github.com/ThuanyMendonca/project/model"
	"github.com/ThuanyMendonca/project/utils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
)

type BalanceControllerSpy struct {
	business balance.IBalanceBusiness
}

func TestPost(t *testing.T) {
	validPaylod := &model.Balance{}
	gofakeit.Struct(validPaylod)

	balance.Router(&gin.Default().RouterGroup)
	tests := []struct {
		description        string
		payload            *bytes.Buffer
		balance            *model.Balance
		codeStatusExpected int
		errorExpected      string
		services           BalanceControllerSpy
	}{
		{
			description:        "Should return an error if sent invalid payload",
			payload:            utils.MockInvalidPayload(),
			codeStatusExpected: http.StatusBadRequest,
			errorExpected:      "conteúdo da requisição inválido",
		},
		{
			description: "Should return body invalid amount",
			payload: utils.MockValidPayload(&model.Balance{
				Amount: -3,
			}),
			codeStatusExpected: http.StatusBadRequest,
			errorExpected:      "valor é obrigatório",
		},
		{
			description:        "Should return error to create balance",
			payload:            utils.MockValidPayload(validPaylod),
			codeStatusExpected: http.StatusInternalServerError,
			errorExpected:      "ocorreu um erro para inserir o registro",
			services: BalanceControllerSpy{
				business: &balance.BalanceBusinessSpy{
					CreateStatus: http.StatusInternalServerError,
					CreateErr:    errors.New("ocorreu um erro para inserir o registro"),
				},
			},
		},
		{
			description:        "Should return success to create balance",
			payload:            utils.MockValidPayload(validPaylod),
			codeStatusExpected: http.StatusOK,
			services: BalanceControllerSpy{
				business: &balance.BalanceBusinessSpy{
					CreateStatus: http.StatusOK,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Post - %s", test.description), func(t *testing.T) {
			c := balance.NewBalanceController(test.services.business)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = httptest.NewRequest(http.MethodPost, "/balance", test.payload)

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

func TestUpdateController(t *testing.T) {
	validBody := &model.BalanceAmount{}
	gofakeit.Struct(validBody)

	balance.Router(&gin.Default().RouterGroup)
	tests := []struct {
		description        string
		userId             string
		key                string
		payload            *bytes.Buffer
		codeStatusExpected int
		errorExpected      string
		services           BalanceControllerSpy
	}{
		{
			description:        "Should return an error if not sent userId",
			codeStatusExpected: http.StatusBadRequest,
			errorExpected:      "id do usuário é obrigatório",
			payload:            utils.MockValidPayload(validBody),
		},
		{
			description:        "Should return an error if not sent payload",
			userId:             "123",
			key:                "userId",
			codeStatusExpected: http.StatusBadRequest,
			errorExpected:      "conteúdo da requisição inválido",
			payload:            utils.MockInvalidPayload(),
		},
		{
			description:        "Should return an error to update balance",
			userId:             "123",
			key:                "userId",
			codeStatusExpected: http.StatusInternalServerError,
			errorExpected:      "ocorreu um erro para atualizar",
			payload:            utils.MockValidPayload(validBody),
			services: BalanceControllerSpy{
				business: &balance.BalanceBusinessSpy{
					UpdateStatusCode: http.StatusInternalServerError,
					UpdateErr:        errors.New("ocorreu um erro para atualizar"),
				},
			},
		},
		{
			description:        "Should return success to update balance",
			userId:             "123",
			key:                "userId",
			codeStatusExpected: http.StatusOK,
			payload:            utils.MockValidPayload(validBody),
			services: BalanceControllerSpy{
				business: &balance.BalanceBusinessSpy{
					UpdateStatusCode: http.StatusOK,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Post - %s", test.description), func(t *testing.T) {
			c := balance.NewBalanceController(test.services.business)

			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = append(ctx.Params, gin.Param{
				Key:   test.key,
				Value: test.userId,
			})

			ctx.Request = httptest.NewRequest(http.MethodPut, "/balance/", test.payload)

			c.Update(ctx)

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

func TestGetController(t *testing.T) {
	validRep := &model.BalanceResp{}
	gofakeit.Struct(validRep)

	balance.Router(&gin.Default().RouterGroup)
	tests := []struct {
		description        string
		userId             string
		key                string
		payload            *bytes.Buffer
		codeStatusExpected int
		errorExpected      string
		services           BalanceControllerSpy
	}{
		{
			description:        "Should return an error if not sent userId",
			codeStatusExpected: http.StatusBadRequest,
			errorExpected:      "id do usuário é obrigatório",
			payload:            utils.MockValidPayload(validRep),
		},
		{
			description:        "Should return an error to get balance",
			userId:             "123",
			key:                "userId",
			codeStatusExpected: http.StatusInternalServerError,
			errorExpected:      "ocorreu um erro para pegar o registro",
			payload:            utils.MockValidPayload(validRep),
			services: BalanceControllerSpy{
				business: &balance.BalanceBusinessSpy{
					GetStatusCode:  http.StatusInternalServerError,
					GetError:       errors.New("ocorreu um erro para pegar o registro"),
					GetBalanceResp: nil,
				},
			},
		},
		{
			description:        "Should return success to get balance",
			userId:             "123",
			key:                "userId",
			codeStatusExpected: http.StatusOK,
			payload:            utils.MockValidPayload(validRep),
			services: BalanceControllerSpy{
				business: &balance.BalanceBusinessSpy{
					GetStatusCode:  http.StatusOK,
					GetBalanceResp: validRep,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Post - %s", test.description), func(t *testing.T) {
			c := balance.NewBalanceController(test.services.business)

			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = append(ctx.Params, gin.Param{
				Key:   test.key,
				Value: test.userId,
			})

			ctx.Request = httptest.NewRequest(http.MethodPut, "/balance/", test.payload)

			c.Get(ctx)

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
