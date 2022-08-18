package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ThuanyMendonca/project/model"
	"github.com/gin-gonic/gin"
)

type IUserController interface {
	Post(c *gin.Context)
}

type UserController struct {
	userBusiness IUserBusiness
}

func NewUserController(userBusiness IUserBusiness) IUserController {
	return &UserController{userBusiness}
}

func (u *UserController) Post(c *gin.Context) {
	user := &model.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("conteúdo da requisição inválido"))
		return
	}

	if err := user.IsValid(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := u.userBusiness.Post(user)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	if statusCode != http.StatusCreated {
		c.AbortWithStatus(statusCode)
		return
	}
	c.AbortWithStatus(statusCode)

}

func (u *UserController) Inactive(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("id é obrigatório"))
		return
	}

	id64, _ := strconv.ParseInt(id, 10, 64)

	statusCode, err := u.userBusiness.Inactive(id64)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, err.Error())
		return
	}

	c.AbortWithStatus(statusCode)

}
