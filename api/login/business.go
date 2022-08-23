package login

import (
	"errors"
	"net/http"

	"github.com/ThuanyMendonca/project/model"
	"github.com/ThuanyMendonca/project/repository"
	"github.com/ThuanyMendonca/project/service/jwtAuth"
	"github.com/ThuanyMendonca/project/utils"
	"gorm.io/gorm"
)

type ILoginBusiness interface {
	Login(login *model.Login) (int, *model.TokenResp, error)
}

type LoginBusiness struct {
	userRepo repository.IUserRepository
}

func NewLoginBusiness(userRepo repository.IUserRepository) ILoginBusiness {
	return &LoginBusiness{userRepo}
}

func (l *LoginBusiness) Login(login *model.Login) (int, *model.TokenResp, error) {
	user, err := l.userRepo.Get(login.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusInternalServerError, nil, err
		}
		return http.StatusInternalServerError, nil, err
	}

	if err := utils.CompareHash(user.Password, login.Password); err != nil {
		return http.StatusBadRequest, nil, errors.New("credenciais inválidas")
	}

	token, err := jwtAuth.NewJWTService().GenerateToken(uint(user.Id))
	if err != nil {
		return http.StatusInternalServerError, nil, errors.New("ocorreu um erro para gerar o token")
	}
	tokenResp := &model.TokenResp{
		AccessToken: token,
	}
	return http.StatusOK, tokenResp, nil
}
