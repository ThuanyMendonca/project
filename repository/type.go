package repository

import (
	"github.com/ThuanyMendonca/project/model"
	"gorm.io/gorm"
)

type ITypeRepository interface {
	Get(id int64) (*model.Type, error)
}

type TypeRepository struct {
	db *gorm.DB
}

func NewTypeRepository(db *gorm.DB) ITypeRepository {
	return &TypeRepository{db}
}

func (t *TypeRepository) Get(id int64) (*model.Type, error) {
	userType := &model.Type{}
	find := t.db

	if id != 0 {
		find.Where("id = ?", id)
	}

	err := find.Find(userType).Error

	if err != nil {
		return nil, err
	}

	return userType, nil

}
