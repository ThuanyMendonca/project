package user_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	userApi "github.com/ThuanyMendonca/project/api/user"
	"github.com/ThuanyMendonca/project/model"
	"github.com/ThuanyMendonca/project/repository"
	userRepository "github.com/ThuanyMendonca/project/repository"
	"github.com/ThuanyMendonca/project/utils"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

type UserBusinessRepositorySpy struct {
	UserRepository userRepository.IUserRepository
}

func TestPost(t *testing.T) {
	user := &model.User{}
	gofakeit.Struct(user)

	tests := []struct {
		description        string
		userModel          *model.User
		statusCodeExpected int
		erroExpected       error
		repositories       UserBusinessRepositorySpy
		utils              utils.GenerateHashSpy
	}{
		{
			description:        "Should return error to get validation",
			statusCodeExpected: http.StatusInternalServerError,
			userModel:          user,
			erroExpected:       errors.New("falha para validar"),
			repositories: UserBusinessRepositorySpy{
				UserRepository: &repository.UserRepositorySpy{
					GetValidationErr: errors.New("falha para validar"),
				},
			},
		},
		{
			description:        "Should return error when user email or document is exist",
			statusCodeExpected: http.StatusBadRequest,
			userModel:          user,
			erroExpected:       errors.New("não é possível cadastrar usuário com os dados informados, verifique se já possuí cadastro"),
			repositories: UserBusinessRepositorySpy{
				UserRepository: &userRepository.UserRepositorySpy{
					GetValidationResp: user,
				},
			},
		},
		{
			description:        "Should return error to create user",
			statusCodeExpected: http.StatusInternalServerError,
			userModel:          user,
			erroExpected:       errors.New("falha para criar o usuário"),
			repositories: UserBusinessRepositorySpy{
				UserRepository: &userRepository.UserRepositorySpy{
					GetValidationResp: &model.User{Id: 0},
					GetValidationErr:  nil,
					CreateErr:         errors.New("falha para criar o usuário"),
				},
			},
		},
		{
			description:        "Should return success to create user",
			statusCodeExpected: http.StatusCreated,
			userModel:          user,
			erroExpected:       nil,
			repositories: UserBusinessRepositorySpy{
				UserRepository: &userRepository.UserRepositorySpy{
					GetValidationResp: &model.User{Id: 0},
					GetValidationErr:  nil,
					CreateErr:         nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Post - %s", test.description), func(t *testing.T) {
			u := userApi.NewUserBusiness(test.repositories.UserRepository)

			statusCode, err := u.Post(test.userModel)

			assert.Equal(t, test.statusCodeExpected, statusCode)
			assert.Equal(t, test.erroExpected, err)
		})
	}
}

func TestInactive(t *testing.T) {
	tests := []struct {
		description        string
		userId             int64
		statusCodeExpected int
		erroExpected       error
		repositories       UserBusinessRepositorySpy
	}{
		{
			description:        "Should return error to inactive user",
			statusCodeExpected: http.StatusInternalServerError,
			userId:             1,
			erroExpected:       errors.New("falha para inativar o usuário"),
			repositories: UserBusinessRepositorySpy{
				UserRepository: &repository.UserRepositorySpy{
					InactiveErr: errors.New("falha para inativar o usuário"),
				},
			},
		},
		{
			description:        "Should return success to inactive user",
			statusCodeExpected: http.StatusOK,
			userId:             12,
			repositories: UserBusinessRepositorySpy{
				UserRepository: &userRepository.UserRepositorySpy{
					InactiveErr: nil,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Inactive - %s", test.description), func(t *testing.T) {
			u := userApi.NewUserBusiness(test.repositories.UserRepository)

			statusCode, err := u.Inactive(test.userId)

			assert.Equal(t, test.statusCodeExpected, statusCode)
			assert.Equal(t, test.erroExpected, err)
		})
	}
}
