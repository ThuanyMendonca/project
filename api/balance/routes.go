package balance

import (
	"github.com/ThuanyMendonca/project/dependency"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	b := NewBalanceBusiness(dependency.BalanceRepository)
	c := NewBalanceController(b)

	r.POST("", c.Post)
	r.PUT("/:userId", c.Update)
	r.GET("/:userId", c.Get)
}
