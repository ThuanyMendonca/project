package repository

import "github.com/ThuanyMendonca/project/model"

type UserRepositorySpy struct {
	IUserRepository
	CreateErr         error
	GetByIdResp       *model.User
	GetByIdErr        error
	GetValidationResp *model.User
	GetValidationErr  error
	InactiveErr       error
	IsActiveResp      bool
	IsActiveErr       error
	GetTypeResp       string
	GetTypeErr        error
}

func (u *UserRepositorySpy) Create(user *model.User) error {
	return u.CreateErr
}
func (u *UserRepositorySpy) GetById(userId int64) (*model.User, error) {
	return u.GetByIdResp, u.GetByIdErr
}

func (u *UserRepositorySpy) GetValidation(filter *model.User) (*model.User, error) {
	return u.GetValidationResp, u.GetValidationErr
}

func (u *UserRepositorySpy) Inactive(id int64) error {
	return u.InactiveErr
}

func (u *UserRepositorySpy) IsActive(id int64) (bool, error) {
	return u.IsActiveResp, u.IsActiveErr
}

func (u *UserRepositorySpy) GetType(id int64) (string, error) {
	return u.GetTypeResp, u.GetTypeErr
}
