package repository

import (
	"github.com/ThuanyMendonca/project/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *model.User) error
	GetValidation(document, email string) (*model.User, error)
	Inactive(id int64) error
	IsActive(id int64) (bool, error)
	GetType(id int64) (*string, error)
}

type UserRepository struct {
	db       *gorm.DB
	userType TypeRepository
}

func NewUserRepository(db *gorm.DB, userType TypeRepository) IUserRepository {
	return &UserRepository{db, userType}
}

func (u *UserRepository) Create(user *model.User) error {
	return u.db.Create(user).Error
}

func (u *UserRepository) GetValidation(document, email string) (*model.User, error) {
	user := &model.User{}
	find := u.db

	if document != "" || email != "" {
		find = find.Where("document = ? or email = ?", document, email)
	}

	err := find.Find(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) Inactive(id int64) error {
	return u.db.Model(&model.User{}).Where("id = ?", id).Update("is_active", false).Error
}

func (u *UserRepository) IsActive(id int64) (bool, error) {
	user := &model.User{}
	find := u.db

	if id != 0 {
		find = find.Where("id = ? and is_active = ?", id, true)
	}

	err := find.Find(user).Error

	if err != nil || user == nil || (user != nil && user.Id == 0) {
		return false, err
	}

	return true, nil
}

func (u *UserRepository) GetType(id int64) (*string, error) {
	userType := &model.Type{}
	find := u.db

	if id != 0 {
		find = find.Where("id = ?", id)
	}

	find = find.Joins("Type")

	err := find.Find(userType).Error

	if err != nil {
		return nil, err
	}

	return &userType.Description, nil

}
