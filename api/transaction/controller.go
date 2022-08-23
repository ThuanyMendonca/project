package transaction

import (
	"net/http"
	"strconv"

	"github.com/ThuanyMendonca/project/model"
	"github.com/gin-gonic/gin"
)

type ITransactionController interface {
	GetTransactionByUserId(c *gin.Context)
	Post(c *gin.Context)
}

type TransactionController struct {
	transactionBusiness ITransactionBusiness
}

func NewUTransactionController(transactionBusiness ITransactionBusiness) ITransactionController {
	return &TransactionController{transactionBusiness}
}

func (t *TransactionController) GetTransactionByUserId(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "id do usuário é obrigatório"})
		return
	}

	userId64, _ := strconv.ParseInt(userId, 10, 64)

	statusCode, transactions, err := t.transactionBusiness.GetTransactionByUserId(userId64)
	if err != nil {
		c.AbortWithStatusJSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.JSON(statusCode, transactions)
}

func (t *TransactionController) Post(c *gin.Context) {
	transaction := &model.Transaction{}

	if err := c.ShouldBindJSON(transaction); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "conteúdo da requisição inválido"})
		return
	}

	if err := transaction.IsValid(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	statusCode, err := t.transactionBusiness.Create(transaction)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.AbortWithStatusJSON(statusCode, gin.H{"message": "transferência realizada com sucesso"})

}
