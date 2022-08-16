package user

import (
	"errors"
	"net/http"

	"github.com/ThuanyMendonca/project/model"
	userRepository "github.com/ThuanyMendonca/project/repository"
	"gorm.io/gorm"
)

type IUserBusiness interface {
	Post(user *model.User) (int, error)
	Inactive(id int64) (int, error)
}

type UserBusiness struct {
	UserRepository userRepository.IUserRepository
}

func NewUserBusiness(userRepository userRepository.UserRepository) IUserBusiness {
	return &UserBusiness{&userRepository}
}

func (u *UserBusiness) Post(user *model.User) (int, error) {
	verifyUser, err := u.UserRepository.GetValidation(user.Document, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, err
	}

	if verifyUser != nil {
		return http.StatusBadRequest, errors.New("não é possível cadastrar usuário com os dados informados, verifique se já possuí cadastro")
	}

	if err := u.UserRepository.Create(user); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (u *UserBusiness) Inactive(id int64) (int, error) {
	if err := u.UserRepository.Inactive(id); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
