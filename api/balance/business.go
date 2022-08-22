package balance

import (
	"errors"
	"net/http"
	"time"

	"github.com/ThuanyMendonca/project/model"
	"github.com/ThuanyMendonca/project/repository"
	"gorm.io/gorm"
)

type IBalanceBusiness interface {
	Create(userId int64, value float64) (int, error)
	Get(userId int64) (int, *model.BalanceResp, error)
	Update(userId int64, value float64) (int, error)
}

type BalanceBusiness struct {
	balanceRepo repository.IBalanceRepository
}

func NewBalanceBusiness(balanceRepo repository.IBalanceRepository) IBalanceBusiness {
	return &BalanceBusiness{balanceRepo}
}

func (b *BalanceBusiness) Create(userId int64, value float64) (int, error) {
	balance, err := b.balanceRepo.Get(userId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusInternalServerError, err
	}
	if balance != nil {
		return http.StatusBadRequest, errors.New("já existe um registro de saldo, considere atualizar o saldo")
	}

	if err := b.balanceRepo.Add(&model.Balance{
		UserId:    userId,
		Amount:    value,
		CreatedAt: time.Now(),
	}); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil

}

func (b *BalanceBusiness) Update(userId int64, value float64) (int, error) {
	balance, err := b.balanceRepo.Get(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, errors.New("nenhum registro encontrado para ser atualizado")
		}
		return http.StatusInternalServerError, err
	}

	newBalance, err := b.calculateNewValue(balance.Amount, value)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := b.balanceRepo.Update(userId, *newBalance); err != nil {
		return http.StatusInternalServerError, errors.New("ocorreu um erro ao atualizar o saldo")
	}

	return http.StatusOK, nil
}

func (b *BalanceBusiness) calculateNewValue(oldBalance, newBalance float64) (*float64, error) {
	if newBalance <= 0 {
		return nil, errors.New("não é possível adicionar saldo com valor zerado")
	}

	calculate := oldBalance + newBalance

	return &calculate, nil

}

// Fazer get balance
func (b *BalanceBusiness) Get(userId int64) (int, *model.BalanceResp, error) {
	balance, err := b.balanceRepo.Get(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, nil, nil
		}
		return http.StatusInternalServerError, nil, err
	}

	if balance.Id == 0 {
		return http.StatusNotFound, nil, nil
	}

	balanceResp := &model.BalanceResp{
		UserId: balance.UserId,
		Amount: balance.Amount,
	}
	return http.StatusOK, balanceResp, nil

}
