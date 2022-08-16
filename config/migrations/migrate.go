package migrations

import (
	"github.com/ThuanyMendonca/project/model"
	"gorm.io/gorm"
)

func Load(db *gorm.DB) {
	migrate(db)
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
	)
}
