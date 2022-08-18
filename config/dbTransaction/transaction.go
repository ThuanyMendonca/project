package dbTransaction

import "gorm.io/gorm"

type IDbTransaction interface {
	Begin() (*gorm.DB, error)
	Commit(db *gorm.DB)
	Rollback(db *gorm.DB)
}

type DbTransaction struct {
	db *gorm.DB
}

func NewDbTransaction(db *gorm.DB) IDbTransaction {
	return &DbTransaction{db}
}

func (d *DbTransaction) Begin() (*gorm.DB, error) {
	tx := d.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tx, nil
}

func (d *DbTransaction) Commit(db *gorm.DB) {
	db.Commit()
}

func (d *DbTransaction) Rollback(db *gorm.DB) {
	db.Rollback()
}
