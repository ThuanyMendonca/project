package balance

import (
	"net/http"
	"strconv"

	"github.com/ThuanyMendonca/project/model"
	"github.com/gin-gonic/gin"
)

type IBalanceController interface {
	Post(c *gin.Context)
	Update(c *gin.Context)
	Get(c *gin.Context)
}

type BalanceController struct {
	balanceBusiness IBalanceBusiness
}

func NewBalanceController(balanceBusiness IBalanceBusiness) IBalanceController {
	return &BalanceController{balanceBusiness}
}

func (b *BalanceController) Post(c *gin.Context) {
	balance := &model.Balance{}

	if err := c.ShouldBindJSON(balance); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "conteúdo da requisição inválido"})
		return
	}

	if err := balance.IsValid(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	statusCode, err := b.balanceBusiness.Create(balance.UserId, balance.Amount)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.AbortWithStatus(statusCode)
}

func (b *BalanceController) Update(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "id do usuário é obrigatório"})
		return
	}

	userId64, _ := strconv.ParseInt(userId, 10, 64)

	balance := &model.BalanceAmount{}

	if err := c.ShouldBindJSON(balance); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "conteúdo da requisição inválido"})
		return
	}

	statusCode, err := b.balanceBusiness.Update(userId64, balance.Amount)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.AbortWithStatus(statusCode)
}

func (b *BalanceController) Get(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "id do usuário é obrigatório"})
		return
	}

	userId64, _ := strconv.ParseInt(userId, 10, 64)

	statusCode, balance, err := b.balanceBusiness.Get(userId64)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.JSON(statusCode, balance)

}
