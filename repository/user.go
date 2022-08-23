package repository

import (
	"github.com/ThuanyMendonca/project/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *model.User) error
	GetById(userId int64) (*model.User, error)
	GetValidation(filter *model.User) (*model.User, error)
	Inactive(id int64) error
	IsActive(id int64) (bool, error)
	GetType(id int64) (string, error)
	Get(username string) (*model.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) Create(user *model.User) error {
	return u.db.Create(user).Error
}

func (u *UserRepository) GetValidation(filter *model.User) (*model.User, error) {
	user := &model.User{}
	find := u.db

	if filter.Document != "" || filter.Email != "" {
		find = find.Where("document = ?", filter.Document).Or("email = ?", filter.Email)
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

func (u *UserRepository) GetById(userId int64) (*model.User, error) {
	user := &model.User{}
	find := u.db

	if userId != 0 {
		find = find.Where("id = ?", userId)
	}

	err := find.Find(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) GetType(userTypeId int64) (string, error) {
	userType := &model.Type{}
	find := u.db

	if userTypeId != 0 {
		find = find.Where("id = ?", userTypeId)
	}

	err := find.Find(userType).Error

	if err != nil {
		return "", err
	}

	return userType.Description, nil

}

func (u *UserRepository) Get(username string) (*model.User, error) {
	user := &model.User{}
	find := u.db

	if username != "" {
		find = find.Where("name = ?", username)
	}

	err := find.Find(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
