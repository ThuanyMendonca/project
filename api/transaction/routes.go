package transaction

import (
	"github.com/ThuanyMendonca/project/dependency"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	b := NewTransactionBusiness(dependency.DbTransaction, dependency.TransactionRepository, dependency.UserRepository, dependency.TypeRepository, dependency.BalanceRepository, dependency.AuthorizationService)

	c := NewUTransactionController(b)

	r.POST("", c.Post)
}
