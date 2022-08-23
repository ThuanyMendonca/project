package user

import (
	"errors"
	"net/http"

	"github.com/ThuanyMendonca/project/model"
	userRepository "github.com/ThuanyMendonca/project/repository"
	"github.com/ThuanyMendonca/project/utils"
	"gorm.io/gorm"
)

type IUserBusiness interface {
	Post(user *model.User) (int, error)
	Inactive(id int64) (int, error)
}

type UserBusiness struct {
	userRepository userRepository.IUserRepository
}

func NewUserBusiness(userRepository userRepository.IUserRepository) IUserBusiness {
	return &UserBusiness{userRepository}
}

func (u *UserBusiness) Post(user *model.User) (int, error) {
	verifyUser, err := u.userRepository.GetValidation(&model.User{
		Email:    user.Email,
		Document: user.Document,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, err
	}

	if verifyUser != nil && verifyUser.Id != 0 {
		return http.StatusBadRequest, errors.New("não é possível cadastrar usuário com os dados informados, verifique se já possuí cadastro")
	}

	newPassword, err := utils.GenerateHash(user.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	user.Password = string(newPassword)

	if err := u.userRepository.Create(user); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (u *UserBusiness) Inactive(id int64) (int, error) {
	if err := u.userRepository.Inactive(id); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
