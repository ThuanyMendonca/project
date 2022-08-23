package user

import "github.com/ThuanyMendonca/project/model"

type UserBusinessSpy struct {
	IUserBusiness
	PostStatusCode     int
	PostErr            error
	InactiveStatusCode int
	InactiveErr        error
}

func (u *UserBusinessSpy) Post(user *model.User) (int, error) {
	return u.PostStatusCode, u.PostErr
}

func (u *UserBusinessSpy) Inactive(id int64) (int, error) {
	return u.InactiveStatusCode, u.InactiveErr
}
