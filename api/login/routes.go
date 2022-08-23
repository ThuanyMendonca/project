package login

import (
	"github.com/ThuanyMendonca/project/dependency"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	b := NewLoginBusiness(dependency.UserRepository)

	c := NewLoginController(b)

	r.POST("", c.Post)
}
