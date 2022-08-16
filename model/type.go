package model

type Type struct {
	Id          int64  `gorm:"column:id;primary_key;" json:"id"`
	Description string `gorm:"column:description;type:varchar(50);not null;" json:"description"`
}
