package model

import (
	"errors"
	"time"
)

type Balance struct {
	Id        int64     `gorm:"column:id;primary_key;" json:"id"`
	UserId    int64     `gorm:"column:user_id;not null"`
	User      User      `gorm:"foreignKey:user_id;" json:"user_id"`
	Amount    float64   `gorm:"column:value;type:decimal(15,2);not null;" json:"value"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;" json:"updatedAt"`
}

type BalanceResp struct {
	UserId int64   `json:"user_id"`
	Amount float64 `json:"value"`
}

type BalanceAmount struct {
	Amount float64 `json:"value"`
}

func (b *Balance) IsValid() error {
	if b.Amount <= 0 {
		return errors.New("valor é obrigatório")
	}

	if b.UserId == 0 {
		return errors.New("usuário é obrigatório")
	}

	return nil
}
