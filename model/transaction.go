package model

import (
	"errors"
	"time"
)

type Transaction struct {
	Id        int64     `gorm:"column:id;primary_key;" json:"id"`
	Value     float64   `gorm:"column:value;type:decimal(15,2);not null;" json:"value"`
	PayerId   int64     `gorm:"column:payer_id;not null"`
	Payer     User      `gorm:"foreignKey:payer_id;" json:"payer"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;" json:"updatedAt"`
}

func (t *Transaction) IsValid() error {
	if t.Value <= 0 {
		return errors.New("valor é obrigatório")
	}

	if t.PayerId == 0 {
		return errors.New("documento é obrigatório")
	}

	return nil
}
