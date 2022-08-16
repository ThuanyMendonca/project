package model

type Balance struct {
	Id     int64   `gorm:"column:id;primary_key;" json:"id"`
	UserId int64   `gorm:"column:user_id;not null"`
	User   User    `gorm:"foreignKey:user_id;" json:"user_id"`
	Amount float64 `gorm:"column:value;type:decimal(15,2);not null;" json:"value"`
}
