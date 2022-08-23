package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThuanyMendonca/project/api/user"
	"github.com/ThuanyMendonca/project/model"
	"github.com/ThuanyMendonca/project/utils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type UserControllerSpy struct {
	userBusiness user.IUserBusiness
}

func TestPostController(t *testing.T) {
	response := &model.User{}
	gofakeit.Slice(response)

	user.Router(&gin.Default().RouterGroup)

	tests := []struct {
		description        string
		payload            *bytes.Buffer
		codeStatusExpected int
		errorExpected      string
		services           UserControllerSpy
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
			errorExpected:      "nome é obrigatório",
			payload: utils.MockValidPayload(&model.User{
				Name: "",
			}),
		},
		{
			description:        "Should return an error to post user",
			codeStatusExpected: http.StatusInternalServerError,
			errorExpected:      "ocorreu um erro para criar o usuário",
			payload:            utils.MockValidPayload(response),
			services: UserControllerSpy{
				userBusiness: &user.UserBusinessSpy{
					PostStatusCode: http.StatusInternalServerError,
					PostErr:        errors.New("ocorreu um erro para criar o usuário"),
				},
			},
		},
		{
			description:        "Should return success to post user",
			codeStatusExpected: http.StatusCreated,
			payload:            utils.MockValidPayload(response),
			services: UserControllerSpy{
				userBusiness: &user.UserBusinessSpy{
					PostStatusCode: http.StatusCreated,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Post - %s", test.description), func(t *testing.T) {
			c := user.NewUserController(test.services.userBusiness)

			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)

			ctx.Request = httptest.NewRequest(http.MethodPost, "/user", test.payload)

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

func TestInactiveController(t *testing.T) {
	user.Router(&gin.Default().RouterGroup)

	tests := []struct {
		description        string
		userId             string
		key                string
		codeStatusExpected int
		errorExpected      string
		services           UserControllerSpy
	}{
		{
			description:        "Should return an error if not sent id",
			codeStatusExpected: http.StatusBadRequest,
			errorExpected:      "id do usuário é obrigatório",
		},
		{
			description:        "Should return an error to inactive user",
			userId:             "123",
			key:                "id",
			codeStatusExpected: http.StatusInternalServerError,
			errorExpected:      "ocorreu um erro para inativar o usuário",
			services: UserControllerSpy{
				userBusiness: &user.UserBusinessSpy{
					InactiveStatusCode: http.StatusInternalServerError,
					InactiveErr:        errors.New("ocorreu um erro para inativar o usuário"),
				},
			},
		},
		{
			description:        "Should return success to get transactions",
			userId:             "123",
			key:                "id",
			codeStatusExpected: http.StatusOK,
			services: UserControllerSpy{
				userBusiness: &user.UserBusinessSpy{
					InactiveStatusCode: http.StatusOK,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Inactive - %s", test.description), func(t *testing.T) {
			c := user.NewUserController(test.services.userBusiness)

			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = append(ctx.Params, gin.Param{
				Key:   test.key,
				Value: test.userId,
			})

			ctx.Request = httptest.NewRequest(http.MethodPut, "/user/", nil)

			c.Inactive(ctx)

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
