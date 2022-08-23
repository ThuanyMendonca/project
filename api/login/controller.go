package login

import (
	"net/http"

	"github.com/ThuanyMendonca/project/model"
	"github.com/gin-gonic/gin"
)

type ILoginController interface {
	Post(c *gin.Context)
}

type LoginController struct {
	loginBusiness ILoginBusiness
}

func NewLoginController(loginBusiness ILoginBusiness) ILoginController {
	return &LoginController{loginBusiness}
}

func (l *LoginController) Post(c *gin.Context) {
	credentials := &model.Login{}

	if err := c.ShouldBindJSON(credentials); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "conteúdo da requisição inválido"})
		return
	}

	statusCode, token, err := l.loginBusiness.Login(credentials)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(statusCode, token)

}
