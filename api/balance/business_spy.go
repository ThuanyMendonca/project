package balance

import "github.com/ThuanyMendonca/project/model"

type BalanceBusinessSpy struct {
	IBalanceBusiness
	CreateStatus     int
	CreateErr        error
	GetStatusCode    int
	GetBalanceResp   *model.BalanceResp
	GetError         error
	UpdateStatusCode int
	UpdateErr        error
}

func (b *BalanceBusinessSpy) Create(userId int64, value float64) (int, error) {
	return b.CreateStatus, b.CreateErr
}

func (b *BalanceBusinessSpy) Get(userId int64) (int, *model.BalanceResp, error) {
	return b.GetStatusCode, b.GetBalanceResp, b.GetError
}

func (b *BalanceBusinessSpy) Update(userId int64, value float64) (int, error) {
	return b.UpdateStatusCode, b.UpdateErr
}
