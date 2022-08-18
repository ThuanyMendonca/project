package transaction

import (
	"errors"
	"net/http"

	"github.com/ThuanyMendonca/project/model"
	"github.com/gin-gonic/gin"
)

type ITransactionController interface {
	Post(c *gin.Context)
}

type TransactionController struct {
	transactionBusiness ITransaction
}

func NewUTransactionController(transactionBusiness ITransaction) ITransactionController {
	return &TransactionController{transactionBusiness}
}

func (t *TransactionController) Post(c *gin.Context) {
	transaction := &model.Transaction{}

	if err := c.ShouldBindJSON(transaction); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("conteúdo da requisição inválido"))
		return
	}

	if err := transaction.IsValid(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	statusCode, err := t.transactionBusiness.Create(transaction)
	if err != nil {
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.AbortWithStatusJSON(statusCode, gin.H{"message": "transferência realizada com sucesso"})

}
