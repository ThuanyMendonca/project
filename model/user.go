package model

import (
	"errors"
	"time"
)

type User struct {
	Id        int64     `gorm:"column:id;primary_key;" json:"id"`
	Name      string    `gorm:"column:name;type:varchar(150);not null;" json:"name"`
	TypeId    int64     `gorm:"column:type_id;not null"`
	Type      Type      `gorm:"foreignKey:type_id;" json:"type"`
	Document  string    `gorm:"column:document;type:varchar(50);not null; unique;" json:"document"`
	Email     string    `gorm:"column:email;type:varchar(100);not null; unique;" json:"email"`
	Password  string    `gorm:"column:password;type:varchar(100);not null;" json:"password"`
	IsActive  bool      `gorm:"column:is_active;default:true;type:bool;not null;" json:"isActive"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null" json:"updatedAt"`
}

func (u *User) IsValid() error {
	if u.Name == "" {
		return errors.New("nome é obrigatório")
	}

	if u.Document == "" {
		return errors.New("documento é obrigatório")
	}

	if u.Email == "" {
		return errors.New("e-mail é obrigatório")
	}

	if u.TypeId == 0 {
		return errors.New("tipo de usuário é obrigatório")
	}

	return nil
}
