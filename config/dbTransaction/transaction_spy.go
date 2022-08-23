package dbTransaction

import "gorm.io/gorm"

type DbTransactionSpy struct {
	BeginResponse *gorm.DB
	BeginErr      error
}

func (d *DbTransactionSpy) Begin() (*gorm.DB, error) {
	return d.BeginResponse, d.BeginErr
}

func (d *DbTransactionSpy) Commit(db *gorm.DB) {
}

func (d *DbTransactionSpy) Rollback(db *gorm.DB) {

}
